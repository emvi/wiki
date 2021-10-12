export const ToastStore = {
    state: {
        toastException: false,
        toastType: "red", // blue, green, orange, red,
        toastMessage: "",
        toastTimeout: null
    },
    mutations: {
        setToastMessage(state, message) {
            state.toastMessage = message;
        },
        setToastException(state, exception) {
            state.toastException = exception;
        },
        setToastType(state, type) {
            state.toastType = type;
        },
        setToastTimeout(state, timeout) {
            state.toastTimeout = timeout;
        }
    },
    actions: {
        showToast(context, {exception, type, message, time}) {
            let toastType = "red";
            let toastTime = 5000;

            if(type) {
                toastType = type;
            }

            if(time) {
                toastTime = time;
            }

            clearTimeout(context.state.toastTimeout);
            context.commit("setToastException", exception);
            context.commit("setToastMessage", message);
            context.commit("setToastType", toastType);
            context.commit("setToastTimeout", setTimeout(() => {
                context.dispatch("resetToast");
            }, toastTime));
        },
        resetToast(context) {
            context.commit("setToastMessage", "");
        }
    },
    getters: {
        toastException(state) {
            return state.toastException;
        },
        toastMessage(state) {
            return state.toastMessage;
        },
        toastType(state) {
            return state.toastType;
        }
    }
};
