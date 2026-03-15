#!/usr/bin/env bash
# seed script - creates test data via the api
# usage: ./scripts/seed.sh

set -eo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
API="http://localhost:8080"

DB_URL=$(grep '^DATABASE_URL=' "$SCRIPT_DIR/../.env" | cut -d'=' -f2-)

if [ -z "$DB_URL" ]; then
  echo "error: DATABASE_URL not found in .env"
  exit 1
fi

echo "==> registering test admin user..."
REG=$(curl -sf "$API/api/auth/register" \
  -H 'Content-Type: application/json' \
  -d '{"email":"admin@retrosnack.shop","password":"testpassword123"}' 2>/dev/null) || {
  echo "    user may already exist, trying login..."
  REG=$(curl -sf "$API/api/auth/login" \
    -H 'Content-Type: application/json' \
    -d '{"email":"admin@retrosnack.shop","password":"testpassword123"}')
}

TOKEN=$(echo "$REG" | python3 -c "import sys,json; print(json.load(sys.stdin)['token'])")
USER_ID=$(echo "$REG" | python3 -c "import sys,json; print(json.load(sys.stdin)['user']['id'])")
echo "    user: $USER_ID"

echo "==> promoting user to admin..."
psql "$DB_URL" -c "UPDATE users SET role = 'admin' WHERE id = '$USER_ID';" 2>/dev/null || {
  echo "    psql failed - run manually:"
  echo "    psql \"$DB_URL\" -c \"UPDATE users SET role = 'admin' WHERE id = '$USER_ID';\""
}

echo "==> re-logging in for admin token..."
REG=$(curl -sf "$API/api/auth/login" \
  -H 'Content-Type: application/json' \
  -d '{"email":"admin@retrosnack.shop","password":"testpassword123"}')
TOKEN=$(echo "$REG" | python3 -c "import sys,json; print(json.load(sys.stdin)['token'])")
AUTH="Authorization: Bearer $TOKEN"

echo "==> ensuring categories exist..."
ensure_category() {
  local name="$1" slug="$2"
  CAT_ID=$(psql "$DB_URL" -t -A -c "
    INSERT INTO categories (id, name, slug) VALUES (gen_random_uuid(), '$name', '$slug')
    ON CONFLICT (slug) DO UPDATE SET name = '$name'
    RETURNING id;
  " | head -1 | tr -d '[:space:]')
  echo "    category: $name ($CAT_ID)"
}

ensure_category "tops" "tops"
TOPS_CAT="$CAT_ID"
ensure_category "bottoms" "bottoms"
BOTTOMS_CAT="$CAT_ID"
ensure_category "sets" "sets"
SETS_CAT="$CAT_ID"
ensure_category "sweaters" "sweaters"
SWEATERS_CAT="$CAT_ID"

echo "==> creating cozy edition drop..."
DROP=$(curl -sf "$API/api/drops" \
  -H 'Content-Type: application/json' \
  -H "$AUTH" \
  -d '{
    "name": "cozy edition",
    "slug": "cozy-edition",
    "description": "warm, soft, and cozy pieces for the colder months.",
    "instagram_url": ""
  }')
DROP_ID=$(echo "$DROP" | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])")
echo "    drop: cozy edition ($DROP_ID)"

create_product() {
  local title="$1" brand="$2" condition="$3" price="$4" cat_id="$5" notes="$6" drop_id="$7"

  local json
  json=$(python3 -c "
import json
d = {
    'title': '$title',
    'brand': '$brand',
    'condition': '$condition',
    'price_cents': $price,
    'category_id': '$cat_id',
    'description': '',
    'instagram_post_url': '',
    'notes': '$notes'
}
if '$drop_id':
    d['drop_id'] = '$drop_id'
print(json.dumps(d))
")

  PROD=$(curl -sf "$API/api/products" \
    -H 'Content-Type: application/json' \
    -H "$AUTH" \
    -d "$json")

  PROD_ID=$(echo "$PROD" | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])")
  echo "    product: $title ($PROD_ID)"
}

create_variant() {
  local prod_id="$1" size="$2" color="$3" sku="$4"

  VAR=$(curl -sf "$API/api/products/$prod_id/variants" \
    -H 'Content-Type: application/json' \
    -H "$AUTH" \
    -d "{\"size\": \"$size\", \"color\": \"$color\", \"sku\": \"$sku\"}")

  VAR_ID=$(echo "$VAR" | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])")

  curl -sf "$API/api/variants/$VAR_ID/stock" \
    -X PUT \
    -H 'Content-Type: application/json' \
    -H "$AUTH" \
    -d '{"quantity": 1}' > /dev/null

  echo "      variant: $size ($VAR_ID)"
}

echo "==> creating products..."

# cozy edition - fleece sweater
create_product "fleece sweater" "" "good" 1800 "$SWEATERS_CAT" "" "$DROP_ID"
create_variant "$PROD_ID" "M" "" "fleece-sweater-m"

# cozy edition - brandy cargos
create_product "brandy cargos" "brandy melville" "good" 3500 "$BOTTOMS_CAT" "" "$DROP_ID"
create_variant "$PROD_ID" "S" "" "brandy-cargos-s"

# cozy edition - white long sleeve
create_product "white long sleeve" "" "new" 2200 "$TOPS_CAT" "" "$DROP_ID"
create_variant "$PROD_ID" "M" "" "white-long-sleeve-m"

# cozy edition - zara teal sweats
create_product "zara teal sweats" "zara" "good" 2000 "$BOTTOMS_CAT" "" "$DROP_ID"
create_variant "$PROD_ID" "S" "" "zara-teal-sweats-s"

# cozy edition - lulu wool sweater
create_product "lulu wool sweater" "lululemon" "new" 6500 "$SWEATERS_CAT" "" "$DROP_ID"
create_variant "$PROD_ID" "M/L" "" "lulu-wool-sweater-ml"

# cozy edition - lulu gym set
create_product "lulu gym set" "lululemon" "new" 6500 "$SETS_CAT" "" "$DROP_ID"
create_variant "$PROD_ID" "size 4 bra" "" "lulu-gym-set-bra-4"
create_variant "$PROD_ID" "size 6 leggings" "" "lulu-gym-set-leg-6"

# standalone - lace top
create_product "lace top" "" "good" 2500 "$TOPS_CAT" "has padding" ""
create_variant "$PROD_ID" "XS" "" "lace-top-xs"

echo ""
echo "==> done! 7 products seeded."
echo "    admin token: $TOKEN"
echo ""
echo "    test the api:"
echo "    curl '$API/api/products' | python3 -m json.tool"
echo "    curl '$API/api/drops/cozy-edition/products' | python3 -m json.tool"
