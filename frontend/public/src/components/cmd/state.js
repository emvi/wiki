export const CmdStore = {
    state: {
        cmdOpen: false,
        closeCmd: false,
        cmd: "",
        typeahead: "",
        error: null,
        column: [] // {view, row, meta}, ...
    },
    mutations: {
        setCmdOpen(state, open) {
            state.cmdOpen = open;
        },
        setCloseCmd(state, close) {
            state.closeCmd = close;
        },
        setCmd(state, cmd) {
            state.cmd = cmd;
        },
        setTypeahead(state, typeahead) {
            state.typeahead = typeahead;
        },
        setRow(state, index) {
            if(state.column.length) {
                state.column[state.column.length-1].row = index;
            }
        },
        setCmdMeta(state, {key, value}) {
            if(state.column.length) {
                state.column[state.column.length-1].meta.set(key, value);
            }
        },
        setCmdMetaVars(state, array) {
            if(state.column.length) {
                for (let i = 0; i < array.length; i++) {
                    state.column[state.column.length-1].meta.set(array[i].key, array[i].value);
                }
            }
        },
        pushColumn(state, data) {
            state.column.push(data);
        },
        popColumn(state) {
            state.column.pop();
            state.error = null;
        },
        setColumn(state, column) {
            state.column = column;
        },
        setError(state, error) {
            state.error = error || null;
        }
    },
    actions: {
        setCmd(context, cmd) {
            context.commit("setCmd", cmd || "");
        },
        resetCmd(context, cmd) {
            context.commit("setCmd", cmd || "");
            context.commit("setTypeahead", "");
            context.commit("setColumn", []);
            context.commit("setError");
        },
        closeCmd(context) {
            context.commit("setCloseCmd", true);
        },
        typeahead(context, typeahead) {
            context.commit("setTypeahead", typeahead || "");
        },
        selectPreviousRow(context) {
            context.commit("setRow", context.getters.row-1);
        },
        selectNextRow(context) {
            context.commit("setRow", context.getters.row+1);
        },
        selectRow(context, row) {
            if(row < 0) {
                row = 0;
            }

            context.commit("setRow", row);
        },
        setCmdMeta(context, keyValue) {
            context.commit("setCmdMeta", keyValue);
        },
        setCmdMetaVars(context, keyValueArray) {
            context.commit("setCmdMetaVars", keyValueArray);
        },
        pushColumn(context, view) {
            context.commit("pushColumn", {
                view,
                row: 0,
                meta: new Map()
            });
        },
        popColumn(context) {
            if(context.state.column.length <= 1) {
                context.dispatch("resetCmd");
            }
            else {
                context.commit("popColumn");
            }
        },
        setError(context, error) {
            context.commit("setError", error);
        }
    },
    getters: {
        cmdOpen(state) {
            return state.cmdOpen;
        },
        closeCmd(state) {
            return state.closeCmd;
        },
        cmd(state) {
            return state.cmd;
        },
        typeahead(state) {
            return state.typeahead;
        },
        view(state) {
            if(!state.column.length) {
                return null;
            }

            return state.column[state.column.length-1].view;
        },
        row(state) {
            if(!state.column.length) {
                return 0;
            }

            return state.column[state.column.length-1].row;
        },
        cmdMeta(state) {
            if(!state.column.length) {
                return {};
            }

            return state.column[state.column.length-1].meta;
        },
        columns(state) {
            return state.column.length;
        },
        cmdError(state) {
            return state.error;
        }
    }
};
