export const RecaptchaStore = {
    state: {
        grecaptcha: null
    },
    mutations: {
        setRecaptcha(state, grecaptcha) {
            if(grecaptcha) {
                state.grecaptcha = grecaptcha;
            }
        }
    }
};
