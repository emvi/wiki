export const PageStore = {
    state: {
        selection: 0,

        // The map is used to communicate between the page and the command menu.
        // metaUpdate gets counted up each time the map is changed, so that the page can listen for changes
        // and update in case it's greater one (or in fact, how many times meta data got updated + 1, in case the
        // page sets meta data itself).
        meta: new Map(),
        metaUpdate: 0
    },
    mutations: {
        setSelection(state, index) {
            state.selection = index;
        },
        setMeta(state, {key, value}) {
            state.meta.set(key, value);
            state.metaUpdate++;
        },
        setMetaVars(state, array) {
            for(let i = 0; i < array.length; i++) {
                state.meta.set(array[i].key, array[i].value);
            }

            state.metaUpdate++;
        },
        resetMeta(state) {
            state.meta = new Map();
            state.metaUpdate = 0;
        }
    },
    actions: {
        selectPrevious(context, maxIndex) {
            let selection = context.getters.selection-1;

            if(selection < 0) {
                selection = maxIndex;
            }

            context.commit("setSelection", selection);
        },
        selectNext(context, maxIndex) {
            let selection = context.getters.selection+1;

            if(selection > maxIndex) {
                selection = 0;
            }

            context.commit("setSelection", selection);
        },
        select(context, selection) {
            context.commit("setSelection", selection);
        },
        setMeta(context, keyValue) {
            context.commit("setMeta", keyValue);
        },
        setMetaVars(context, keyValueArray) {
            context.commit("setMetaVars", keyValueArray);
        },
        resetMeta(context) {
            context.commit("resetMeta");
        }
    },
    getters: {
        selection(state) {
            return state.selection;
        },
        metaUpdate(state) {
            return state.metaUpdate;
        }
    }
};
