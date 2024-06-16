import { App } from 'vue';

export default (app: App<Element>) => {
  app.directive('preReClick', {
    mounted(el, binding) {
      el.addEventListener('click', () => {
        if (!el.disabled) {
          el.disabled = true;
          el.loading = true;
          setTimeout(() => {
            el.disabled = false;
          }, binding.value || 2000);
        }
      });
    }
  });
};