let message = $state('');
let visible = $state(false);
let timeout: ReturnType<typeof setTimeout>;

export const toast = {
  get message() {
    return message;
  },

  get visible() {
    return visible;
  },

  show(msg: string, duration = 3000) {
    clearTimeout(timeout);
    message = msg;
    visible = true;
    timeout = setTimeout(() => {
      visible = false;
    }, duration);
  },
};
