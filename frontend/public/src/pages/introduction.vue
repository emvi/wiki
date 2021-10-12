<template>
    <emvi-layout v-on:enter="nextStep" v-on:esc="skip" v-on:left="previousStep" v-on:right="nextStep" v-on:up="scroll(-1)" v-on:down="scroll(1)">
        <emvi-layout-narrow>
            <small>{{$t("title")}} — {{$t("step")}} {{step}} / 4</small>
            <div class="article-content text">
                <template v-if="step === 1">
                    <h1>{{$t("title1")}}</h1>
                    <div class="text" v-html="$t('step1')"></div>
                </template>
                <template v-if="step === 2">
                    <h1>{{$t("title2")}}</h1>
                    <div class="text" v-html="$t('step2')"></div>
                </template>
                <template v-if="step === 3">
                    <h1>{{$t("title3")}}</h1>
                    <div class="text" v-html="$t('step3')"></div>
                </template>
                <template v-if="step === 4">
                    <h1>{{$t("title4")}}</h1>
                    <div class="text" v-html="$t('step4')"></div>
                </template>
            </div>
            <div class="headline" style="margin-top: 24px;">
                <p class="cursor-pointer" v-on:click="skip">
                    {{$t("skip")}}
                    <emvi-shortcut :shortcut="$t('shortcut_skip')" v-on:click="skip"></emvi-shortcut>
                </p>
                <p class="cursor-pointer" v-on:click="nextStep">
                    {{$t("next")}}
                    <emvi-shortcut :shortcut="$t('shortcut_next')" v-on:click="nextStep"></emvi-shortcut>
                </p>
            </div>
        </emvi-layout-narrow>
    </emvi-layout>
</template>

<script>
    import {scroll} from "../util";
    import {UserService} from "../service";
    import {TitleMixin} from "./title";
    import {emviLayout, emviLayoutNarrow, emviShortcut} from "../components";

    const maxStep = 4;

    export default {
        mixins: [TitleMixin],
        components: {emviLayout, emviLayoutNarrow, emviShortcut},
        data() {
            return {
                step: 1
            };
        },
        methods: {
            previousStep() {
                if(this.step > 1) {
                    this.step--;
                }
            },
            nextStep() {
                this.step++;

                if(this.step > maxStep) {
                    this.skip();
                }
            },
            skip() {
                UserService.setIntroduction(false)
                    .then(() => {
                        this.$router.push("/");
                    })
                    .catch(e => {
                        this.showTechnicalError(e);
                    });
            },
            scroll(dir) {
                scroll(dir);
            }
        }
    }
</script>

<style lang="scss" scoped>
    .text {
        margin: 24px 0;
    }
</style>

<i18n>
    {
        "en": {
            "title": "Introduction",
            "skip": "Skip Introduction",
            "next": "Continue",
            "step": "Step",
            "shortcut_skip": "Esc",
            "shortcut_next": "Enter",
            "title1": "Welcome to Emvi",
            "title2": "The Basics",
            "title3": "User Interface and Controls",
            "title4": "Finish Introduction",
            "step1": "<p>Welcome to Emvi, we are happy to have you on board!</p><p>This introduction will guide you through the basic principles of Emvi and explains how to use the user interface.</p><p>You can skip this introduction by pressing <code>Escape</code> or move on to the next step by pressing <code>Enter</code> on your keyboard. You can also use the arrow keys. In case you would like to revisit the introduction later, you can do that from the settings menu.</p>",
            "step2": "<p>Articles (and anything else) can be interconnected by links. To add a link to an article, just type <code>@</code> and search for an object. The object&#39;s title will then be added to the document and kept synchronized in case it gets changed.</p><img src=\"../static/img/mentions.png\" /><p>Tags help you to group articles by a broader topic. In case you would like to create a series of articles, you can use lists to archive that. They can be used to create a step by step tutorial, to name an example.</p><img src=\"../static/img/list.png\" />",
            "step3": "<p>The user interface is best controlled by the keyboard. You will use the central command menu to access actions, navigate the page, and search. You can open it everywhere by pressing <code>Shift</code> + <code>Space</code>.</p><p>Once it&#39;s open, you can navigate the page by typing <code>.</code> followed by the name of the page or use the autocomplete feature by pressing the tab key. When you press <code>Enter</code>, the selected page will be opened. You can switch the selected result using the arrow keys.</p><p>To search, just start typing. You can preview the search results by pressing Tab. Press <code>Enter</code> to open the selected result.</p><img src=\"../static/img/search.svg\" /><p>To run a command, type <code>/</code>. This will list all available commands for the current page. A few of the commands can be accessed from anywhere, like the <code>/logout</code> command. Once you pressed enter, the command will be executed. Most of the commands will open a menu to enter additional information or confirm the action, but there are a few commands that will be executed immediately. In such menus, you can navigate the entries either by using the arrow keys or, if that&#39;s not possible (like in forms), by using the <code>Tab</code> and <code>Shift</code> + <code>Tab</code> keys. <code>Enter</code> usually sends a form or opens a sub-menu. To go back or close the menu, press <code>Escape</code>.</p><p>The entries on the page can be navigated in a similar fashion to the command menu. Once you got used to it, you will be able to control Emvi without even looking. Obviously, Emvi can be controlled by a mouse or on your touch device as well.</p>",
            "step4": "<p>Thanks again for signing up to Emvi! In case you have questions, feedback, or need help, you can contact us through our <code>/support</code> or by <a href=\"mailto:support@emvi.com\" target=\"_blank\" rel=\"noreferrer\">email</a>.</p><p>Press <code>Enter</code>or <code>Escape</code> to finish the introduction. You can then go ahead and start using Emvi and <code>/invite</code> your first members.</p>"
        },
        "de": {
            "title": "Einführung",
            "skip": "Einführung überspringen",
            "next": "Weiter",
            "step": "Schritt",
            "shortcut_skip": "Esc",
            "shortcut_next": "Enter",
            "title1": "Willkommen bei Emvi",
            "title2": "Die Grundlagen",
            "title3": "Benutzeroberfläche und Bedienung",
            "title4": "Einführung abschließen",
            "step1": "<p>Willkommen bei Emvi, wir sind froh dich an Board zu haben!</p><p>Diese Einführung bringt dir die grundlegenden Konzepte von Emvi näher und erklärt, wie die Benutzeroberfläche verwendet wird.</p><p>Du kannst die Einführung mit <code>Escape</code> überspringen oder durch das <code>Enter</code> zum nächsten Schritt fortfahren. Alternativ kannst du auch die Pfeiltasten verwenden. Möchtest du die Einführung später wiederholen wollen, kannst du das über das Einstellungsmenü tun.</p>",
            "step2": "<p>Artikel (und alles andere) können in Artikeln verlinkt werden. Um eine Verlinkung zu einem Artikel hinzufügen, tippe das <code>@</code> Zeichen und suche nach einem Objekt das du verlinken möchtest. Der Titel des Objekts wird dann in das Dokument eingefügt und automatisch synchronisiert, sollte es geändert werden.</p><img src=\"../static/img/mentions.png\" /><p>Tags helfen die Artikel zu größeren Themen zusammenzufassen. Falls du eine Serie von Artikel erstellen möchtest kannst du Listen verwenden. Sie können z.B. verwendet werden um Schritt für Schritt Anleitungen zu erstellen.</p><img src=\"../static/img/list.png\" />",
            "step3": "<p>Die Benutzeroberfläche lässt sich am besten per Tastatur bedienen. Du wirst das zentrale Befehlsmenü verwenden, um auf Aktionen und die Suche zuzugreifen und um zu navigieren. Du kannst es mit Hilfe des Tastenkürzels <code>Shift</code> + <code>Leer</code> jederzeit öffnen.</p><p>Sobald es geöffnet ist kannst du die Seite wechseln indem du einen <code>.</code> eingibst, gefolgt vom Namen der Seite. Oder du verwendest die automatische Vervollständigung indem du <code>Tab</code> drückst. Wenn du <code>Enter</code> drückst, wird die ausgewählte Seite geöffnet. Mit Hilfe der Pfeiltasten kannst du zwischen den Ergebnissen wechseln.</p><p>Für die Suche kannst du einfach anfangen zu tippen. Du kannst dir die Vorschau zu einem Ergebnis mit Hilfe der <code>Tab</code> Taste ansehen. Durch das Drücken von <code>Enter</code> wir das ausgewählte Ergebnis geöffnet.</p><img src=\"../static/img/search.svg\" /><p>Um einen Befehl auszuführen, gibst du ein <code>/</code> ein. Danach werden alle Befehle aufgelistet, die auf der Seite verfügbar sind. Auf einige der Befehle kann von überall aus zugegriffen werden, wie z.B. der <code>/logout</code> Befehl. Sobald du <code>Enter</code> drückst, wird der ausgewählte Befehl ausgeführt. Die meisten der Befehle öffnen ein weiteres Menü in das du zusätzliche Informationen eingeben kannst, oder das dich um die Bestätigung deiner Aktion bittet. Einige der Befehle werden sofort ausgeführt. In den Untermenüs kannst du die Einträge entweder über die Pfeiltasten durchlaufen oder, sollte das nicht möglich sein (wie z.B. in Formularen), über die <code>Tab</code> bzw. <code>Shift</code> + <code>Tab</code> Taste. <code>Enter</code> sendet üblicherweise ein Formular ab oder öffnet ein Untermenü. Um im Befehlsmenü zurück zu gehen, drücke <code>Escape</code>.</p><p>Einträge auf den Seiten können in der selben Weise wie im Befehlsmenü bedient werden. Sobald du die Steuerung verinnerlicht hast, kannst du Emvi ohne hinzusehen bedienen. Natürlich kann Emvi auch mit der Maus und auf Touchgeräten bedient werden.</p>",
            "step4": "<p>Vielen Dank für die Registrierung bei Emvi! Falls du Fragen oder Feedback hast oder auf ein Problem stößt, kannst du uns über den <code>/support</code> Befehl oder per <a href=\"mailto:support@emvi.com\" target=\"_blank\" rel=\"noreferrer\">E-Mail</a> erreichen.</p><p>Drücke Enter oder Escape um die Einführung abzuschließen. Danach kannst du Emvi verwenden und die ersten Mitglieder <code>/einladen</code>.</p>"
        }
    }
</i18n>
