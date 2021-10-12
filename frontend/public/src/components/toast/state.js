const TIMEOUT = 5000;

export const ToastStore = {
    state: {
        toastColor: "",
        toastMessage: "",
        toastTimeout: null
    },
    mutations: {
        setToastMessage(state, {color, message}) {
            this.commit("resetToast");
            state.toastColor = color;
            state.toastMessage = message;
            state.toastTimeout = setTimeout(() => {
                this.commit("resetToast");
            }, TIMEOUT);
        },
        resetToast(state) {
            if(state.toastTimeout !== null) {
                clearTimeout(state.toastTimeout);
                state.toastTimeout = null;
            }
        }
    },
    actions: {
        success(context, message) {
            context.commit("setToastMessage", {color: "green", message});
        },
        error(context, message) {
            context.commit("setToastMessage", {color: "red", message});
        },
        closeToast(context) {
            context.commit("resetToast");
        }
    },
    getters: {
        showToast(state) {
            return state.toastTimeout !== null;
        },
        toastColor(state) {
            return state.toastColor;
        },
        toastMessage(state) {
            return state.toastMessage;
        }
    }
};
