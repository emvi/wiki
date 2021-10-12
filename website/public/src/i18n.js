import VueI18n from "vue-i18n";
import Vue from "vue";
import {getLocale} from "./util/locale";

Vue.use(VueI18n);

const messages = {
    en: {
        button_save: "Save",
        button_discard: "Discard",
        button_cancel: "Cancel",
        button_back: "Back",
        toast_saved: "Saved."
    },
    de: {
        button_save: "Speichern",
        button_discard: "Verwerfen",
        button_cancel: "Abbrechen",
        button_back: "Zurück",
        toast_saved: "Gespeichert."
    }
};

export function getVueI18n() {
    return new VueI18n({
        locale: getLocale(),
        fallbackLocale: "en",
        silentTranslationWarn: true,
        messages
    });
}
