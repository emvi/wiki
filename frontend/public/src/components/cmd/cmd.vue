<template>
    <div class="overlay background" v-if="cmdOpen">
        <div class="cmd">
            <div ref="cmd">
                <div class="command-line">
                    <emvi-cmd-icon v-on:click="back"></emvi-cmd-icon>
                    <emvi-cmd-input v-on:enter="enter"
                                    v-on:tab="tab"
                                    v-on:esc="esc"
                                    v-on:up="up"
                                    v-on:down="down"
                                    v-on:del="del"></emvi-cmd-input>
                    <i class="icon icon-compass action cursor-pointer" v-on:click.stop="navigation" v-if="!view"></i>
                    <i class="icon icon-flash action cursor-pointer" v-on:click.stop="action" v-if="!view"></i>
                    <i class="icon icon-close action cursor-pointer" v-on:click.stop="clear" v-if="view"></i>
                </div>
                <p class="error" v-show="cmdError">
                    {{cmdError}}
                </p>
                <div class="results no-select" v-show="view">
                    <div id="cmd-results" class="results-content">
                        <emvi-cmd-command v-if="view === 'command'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :esc="escEvent"></emvi-cmd-command>
                        <emvi-cmd-navigation v-if="view === 'navigation'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :esc="escEvent"></emvi-cmd-navigation>
                        <emvi-cmd-search v-if="view === 'search'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :esc="escEvent"></emvi-cmd-search>
                        <emvi-cmd-notifications v-if="view === 'notifications'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :esc="escEvent"></emvi-cmd-notifications>
                        <emvi-cmd-settings v-if="view === 'settings'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :esc="escEvent"></emvi-cmd-settings>
                        <emvi-cmd-settings-account v-if="view === 'settings-account'" :esc="escEvent"></emvi-cmd-settings-account>
                        <emvi-cmd-settings-ui v-if="view === 'settings-ui'" :esc="escEvent"></emvi-cmd-settings-ui>
                        <emvi-cmd-settings-leave v-if="view === 'settings-leave'" :esc="escEvent"></emvi-cmd-settings-leave>
                        <emvi-cmd-help v-if="view === 'help'"
                            :up="upEvent"
                            :down="downEvent"
                            :esc="escEvent"></emvi-cmd-help>
                        <emvi-cmd-support v-if="view === 'support'" :esc="escEvent"></emvi-cmd-support>
                        <emvi-cmd-discard v-if="view === 'discard'" :esc="escEvent"></emvi-cmd-discard>
                        <emvi-cmd-publish v-if="view === 'publish'" :esc="escEvent"></emvi-cmd-publish>
                        <emvi-cmd-recommend v-if="view === 'recommend'"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :esc="escEvent"
                            :up="upEvent"
                            :down="downEvent"></emvi-cmd-recommend>
                        <emvi-cmd-invite v-if="view === 'invite'"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :esc="escEvent"
                            :up="upEvent"
                            :down="downEvent"></emvi-cmd-invite>
                        <emvi-cmd-article-permissions v-if="view === 'article-permissions'"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :esc="escEvent"
                            :up="upEvent"
                            :down="downEvent"></emvi-cmd-article-permissions>
                        <emvi-cmd-article-permissions-member v-if="view === 'article-permissions-member'"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :del="delEvent"
                            :esc="escEvent"
                            :up="upEvent"
                            :down="downEvent"></emvi-cmd-article-permissions-member>
                        <emvi-cmd-article-permissions-member-add v-if="view === 'article-permissions-member-add'"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :del="delEvent"
                            :esc="escEvent"
                            :up="upEvent"
                            :down="downEvent"></emvi-cmd-article-permissions-member-add>
                        <emvi-cmd-recommendations v-if="view === 'recommendations'" :esc="escEvent"></emvi-cmd-recommendations>
                        <emvi-cmd-history v-if="view === 'history'"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :del="delEvent"
                            :esc="escEvent"
                            :up="upEvent"
                            :down="downEvent"></emvi-cmd-history>
                        <emvi-cmd-translation v-if="view === 'translation'"
                            :enter="enterEvent"
                            :esc="escEvent"
                            :up="upEvent"
                            :down="downEvent"></emvi-cmd-translation>
                        <emvi-cmd-disconnect v-if="view === 'disconnect'" :esc="escEvent"></emvi-cmd-disconnect>
                        <emvi-cmd-list v-if="view === 'list'" :esc="escEvent"></emvi-cmd-list>
                        <emvi-cmd-list-entries v-if="view === 'list-entries'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :del="delEvent"
                            :esc="escEvent"></emvi-cmd-list-entries>
                        <emvi-cmd-list-entries-add v-if="view === 'list-entries-add'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :del="delEvent"
                            :esc="escEvent"></emvi-cmd-list-entries-add>
                        <emvi-cmd-list-permissions v-if="view === 'list-permissions'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :del="delEvent"
                            :esc="escEvent"></emvi-cmd-list-permissions>
                        <emvi-cmd-list-permissions-add v-if="view === 'list-permissions-add'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :del="delEvent"
                            :esc="escEvent"></emvi-cmd-list-permissions-add>
                        <emvi-cmd-group v-if="view === 'group'" :esc="escEvent"></emvi-cmd-group>
                        <emvi-cmd-group-member v-if="view === 'group-member'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :del="delEvent"
                            :esc="escEvent"></emvi-cmd-group-member>
                        <emvi-cmd-group-member-add v-if="view === 'group-member-add'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :del="delEvent"
                            :esc="escEvent"></emvi-cmd-group-member-add>
                        <emvi-cmd-tag v-if="view === 'tag'" :esc="escEvent"></emvi-cmd-tag>
                        <emvi-cmd-duplicate-article v-if="view === 'duplicate-article'" :esc="escEvent"></emvi-cmd-duplicate-article>
                        <emvi-cmd-export-article v-if="view === 'export-article'" :esc="escEvent"></emvi-cmd-export-article>
                        <emvi-cmd-archive-article v-if="view === 'archive-article'" :esc="escEvent"></emvi-cmd-archive-article>
                        <emvi-cmd-delete-article v-if="view === 'delete-article'" :esc="escEvent"></emvi-cmd-delete-article>
                        <emvi-cmd-delete-list v-if="view === 'delete-list'" :esc="escEvent"></emvi-cmd-delete-list>
                        <emvi-cmd-delete-group v-if="view === 'delete-group'" :esc="escEvent"></emvi-cmd-delete-group>
                        <emvi-cmd-delete-tag v-if="view === 'delete-tag'" :esc="escEvent"></emvi-cmd-delete-tag>
                        <emvi-cmd-administration v-if="view === 'administration'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :esc="escEvent"></emvi-cmd-administration>
                        <emvi-cmd-administration-image v-if="view === 'administration-image'" :esc="escEvent"></emvi-cmd-administration-image>
                        <emvi-cmd-administration-general v-if="view === 'administration-general'" :esc="escEvent"></emvi-cmd-administration-general>
                        <emvi-cmd-administration-permission v-if="view === 'administration-permission'" :esc="escEvent"></emvi-cmd-administration-permission>
                        <emvi-cmd-administration-members v-if="view === 'administration-members'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :del="delEvent"
                            :esc="escEvent"></emvi-cmd-administration-members>
                        <emvi-cmd-administration-members-invitations v-if="view === 'administration-members-invitations'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :esc="escEvent"></emvi-cmd-administration-members-invitations>
                        <emvi-cmd-administration-members-add v-if="view === 'administration-members-add'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :del="delEvent"
                            :esc="escEvent"></emvi-cmd-administration-members-add>
                        <emvi-cmd-administration-invitations v-if="view === 'administration-invitations'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :del="delEvent"
                            :esc="escEvent"></emvi-cmd-administration-invitations>
                        <emvi-cmd-administration-languages v-if="view === 'administration-languages'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :esc="escEvent"></emvi-cmd-administration-languages>
                        <emvi-cmd-administration-languages-add v-if="view === 'administration-languages-add'" :esc="escEvent"></emvi-cmd-administration-languages-add>
                        <emvi-cmd-administration-clients v-if="view === 'administration-clients'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :del="delEvent"
                            :esc="escEvent"></emvi-cmd-administration-clients>
                        <emvi-cmd-administration-clients-add v-if="view === 'administration-clients-add'" :esc="escEvent"></emvi-cmd-administration-clients-add>
                        <emvi-cmd-administration-delete v-if="view === 'administration-delete'" :esc="escEvent"></emvi-cmd-administration-delete>
                        <emvi-cmd-administration-invitation-link v-if="view === 'administration-invitation-link'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :esc="escEvent"></emvi-cmd-administration-invitation-link>
                        <emvi-cmd-bookmarks-articles v-if="view === 'bookmarks-articles'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :del="delEvent"
                            :esc="escEvent"></emvi-cmd-bookmarks-articles>
                        <emvi-cmd-bookmarks-lists v-if="view === 'bookmarks-lists'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :del="delEvent"
                            :esc="escEvent"></emvi-cmd-bookmarks-lists>
                        <emvi-cmd-observed-articles v-if="view === 'observed-articles'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :del="delEvent"
                            :esc="escEvent"></emvi-cmd-observed-articles>
                        <emvi-cmd-observed-lists v-if="view === 'observed-lists'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :del="delEvent"
                            :esc="escEvent"></emvi-cmd-observed-lists>
                        <emvi-cmd-observed-groups v-if="view === 'observed-groups'"
                            :up="upEvent"
                            :down="downEvent"
                            :enter="enterEvent"
                            :del="delEvent"
                            :esc="escEvent"></emvi-cmd-observed-groups>
                        <emvi-cmd-select-user-group v-if="view === 'select-user-group'"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :del="delEvent"
                            :esc="escEvent"
                            :up="upEvent"
                            :down="downEvent"></emvi-cmd-select-user-group>
                        <emvi-cmd-create v-if="view === 'create'"
                            :enter="enterEvent"
                            :tab="tabEvent"
                            :del="delEvent"
                            :esc="escEvent"
                            :up="upEvent"
                            :down="downEvent"></emvi-cmd-create>
                        <emvi-cmd-introduction v-if="view === 'introduction'" :esc="escEvent"></emvi-cmd-introduction>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {getCmdType} from "./type";
    import emviCmdIcon from "./icon.vue";
    import emviCmdInput from "./input.vue";
    import emviCmdCommand from "./command.vue";
    import emviCmdNavigation from "./navigation.vue";
    import emviCmdSearch from "./search.vue";
    import emviCmdDiscard from "./commands/article/discard.vue";
    import emviCmdPublish from "./commands/article/publish.vue";
    import emviCmdRecommend from "./commands/article/recommend.vue";
    import emviCmdInvite from "./commands/article/invite.vue";
    import emviCmdArticlePermissions from "./commands/article/permissions.vue";
    import emviCmdArticlePermissionsMember from "./commands/article/permissions-member.vue";
    import emviCmdArticlePermissionsMemberAdd from "./commands/article/permissions-member-add.vue";
    import emviCmdRecommendations from "./commands/article/recommendations.vue";
    import emviCmdHistory from "./commands/article/history.vue";
    import emviCmdTranslation from "./commands/article/translation.vue";
    import emviCmdDisconnect from "./commands/article/disconnect.vue";
    import emviCmdList from "./commands/list/list.vue";
    import emviCmdListEntries from "./commands/list/entries.vue";
    import emviCmdListEntriesAdd from "./commands/list/entries-add.vue";
    import emviCmdListPermissions from "./commands/list/permissions.vue";
    import emviCmdListPermissionsAdd from "./commands/list/permissions-add.vue";
    import emviCmdGroup from "./commands/group/group.vue";
    import emviCmdGroupMember from "./commands/group/member.vue";
    import emviCmdGroupMemberAdd from "./commands/group/member-add.vue";
    import emviCmdTag from "./commands/tag/tag.vue";
    import emviCmdDuplicateArticle from "./commands/article/duplicate.vue";
    import emviCmdExportArticle from "./commands/article/export.vue";
    import emviCmdArchiveArticle from "./commands/article/archive.vue";
    import emviCmdDeleteArticle from "./commands/article/delete.vue";
    import emviCmdDeleteList from "./commands/list/delete.vue";
    import emviCmdDeleteGroup from "./commands/group/delete.vue";
    import emviCmdDeleteTag from "./commands/tag/delete.vue";
    import emviCmdNotifications from "./commands/notifications.vue";
    import emviCmdSettings from "./commands/settings/settings.vue";
    import emviCmdSettingsAccount from "./commands/settings/account.vue";
    import emviCmdSettingsUi from "./commands/settings/ui.vue";
    import emviCmdSettingsLeave from "./commands/settings/leave.vue";
    import emviCmdHelp from "./commands/help.vue";
    import emviCmdSupport from "./commands/support.vue";
    import emviCmdAdministration from "./commands/administration/administration.vue";
    import emviCmdAdministrationImage from "./commands/administration/image.vue";
    import emviCmdAdministrationGeneral from "./commands/administration/general.vue";
    import emviCmdAdministrationPermission from "./commands/administration/permission.vue";
    import emviCmdAdministrationMembersInvitations from "./commands/administration/members-invitations.vue";
    import emviCmdAdministrationMembers from "./commands/administration/members.vue";
    import emviCmdAdministrationMembersAdd from "./commands/administration/members-add.vue";
    import emviCmdAdministrationInvitations from "./commands/administration/invitations.vue";
    import emviCmdAdministrationLanguages from "./commands/administration/languages.vue";
    import emviCmdAdministrationLanguagesAdd from "./commands/administration/languages-add.vue";
    import emviCmdAdministrationClients from "./commands/administration/clients.vue";
    import emviCmdAdministrationClientsAdd from "./commands/administration/clients-add.vue";
    import emviCmdAdministrationDelete from "./commands/administration/delete.vue";
    import emviCmdAdministrationInvitationLink from "./commands/administration/invitation-link.vue";
    import emviCmdBookmarksArticles from "./commands/bookmarks/articles.vue";
    import emviCmdBookmarksLists from "./commands/bookmarks/lists.vue";
    import emviCmdObservedArticles from "./commands/observed/articles.vue";
    import emviCmdObservedLists from "./commands/observed/lists.vue";
    import emviCmdObservedGroups from "./commands/observed/groups.vue";
    import emviCmdSelectUserGroup from "./commands/select-user-group.vue";
    import emviCmdCreate from "./commands/create.vue";
    import emviCmdIntroduction from "./commands/settings/introduction.vue";

    const cmdShortcutDelay = 200;

    export default {
        components: {
            emviCmdIcon,
            emviCmdInput,
            emviCmdCommand,
            emviCmdNavigation,
            emviCmdSearch,
            emviCmdDiscard,
            emviCmdPublish,
            emviCmdRecommend,
            emviCmdInvite,
            emviCmdArticlePermissions,
            emviCmdArticlePermissionsMember,
            emviCmdArticlePermissionsMemberAdd,
            emviCmdRecommendations,
            emviCmdHistory,
            emviCmdTranslation,
            emviCmdDisconnect,
            emviCmdList,
            emviCmdListEntries,
            emviCmdListEntriesAdd,
            emviCmdListPermissions,
            emviCmdListPermissionsAdd,
            emviCmdGroup,
            emviCmdGroupMember,
            emviCmdGroupMemberAdd,
            emviCmdTag,
            emviCmdDuplicateArticle,
            emviCmdExportArticle,
            emviCmdArchiveArticle,
            emviCmdDeleteArticle,
            emviCmdDeleteList,
            emviCmdDeleteGroup,
            emviCmdDeleteTag,
            emviCmdNotifications,
            emviCmdSettings,
            emviCmdSettingsAccount,
            emviCmdSettingsUi,
            emviCmdSettingsLeave,
            emviCmdHelp,
            emviCmdSupport,
            emviCmdAdministration,
            emviCmdAdministrationImage,
            emviCmdAdministrationGeneral,
            emviCmdAdministrationPermission,
            emviCmdAdministrationMembersInvitations,
            emviCmdAdministrationMembers,
            emviCmdAdministrationMembersAdd,
            emviCmdAdministrationInvitations,
            emviCmdAdministrationLanguages,
            emviCmdAdministrationLanguagesAdd,
            emviCmdAdministrationClients,
            emviCmdAdministrationClientsAdd,
            emviCmdAdministrationDelete,
            emviCmdAdministrationInvitationLink,
            emviCmdBookmarksArticles,
            emviCmdBookmarksLists,
            emviCmdObservedArticles,
            emviCmdObservedLists,
            emviCmdObservedGroups,
            emviCmdSelectUserGroup,
            emviCmdCreate,
            emviCmdIntroduction
        },
        data() {
            return {
                keydownHandler: null,
                clickHandler: null,
                enterEvent: false,
                tabEvent: false,
                escEvent: false,
                upEvent: false,
                downEvent: false,
                delEvent: undefined // carries the keyboard event
            };
        },
        computed: {
            ...mapGetters(["view", "closeCmd", "cmdOpen", "cmd", "columns", "cmdError"])
        },
        watch: {
            closeCmd(close) {
                if(close) {
                    this.close();
                    this.$nextTick(() => {
                        this.$store.commit("setCloseCmd", false);
                    });
                }
            },
            cmd(cmd) {
                // set open in case the command was set programmatically
                if(cmd && !this.cmdOpen) {
                    this.$store.commit("setCmdOpen", true);
                }

                if(this.columns === 0) {
                    let view = getCmdType(cmd);

                    if(view) {
                        this.$store.dispatch("pushColumn", view);
                    }
                }
            }
        },
        mounted() {
            this.bindKeys();
            this.bindMouse();
            this.$store.commit("setCmdOpen", false);
        },
        beforeDestroy() {
            this.unbindKeys();
            this.unbindMouse();
        },
        methods: {
            bindKeys() {
                let lastKeystroke = 0;

                this.keydownHandler = e => {
                    let elapsed = e.timeStamp-lastKeystroke;

                    if(e.code !== "ShiftLeft" && e.code !== "ShiftRight") {
                        lastKeystroke = e.timeStamp;
                    }

                    if(e.shiftKey && e.code === "Space" && elapsed > cmdShortcutDelay) {
                        e.preventDefault();
                        e.stopPropagation();
                        this.$store.commit("setCmdOpen", true);
                    }
                };
                window.addEventListener("keydown", this.keydownHandler);
            },
            bindMouse() {
                this.clickHandler = window.addEventListener("mousedown", e => {
                    if(this.$refs.cmd && e.target !== this.$refs.cmd && !this.$refs.cmd.contains(e.target)) {
                        this.close();
                    }
                });
            },
            unbindKeys() {
                window.removeEventListener("keydown", this.keydownHandler);
            },
            unbindMouse() {
                window.removeEventListener("click", this.clickHandler);
            },
            enter(e) {
                this.enterEvent = e;
                this.$nextTick(() => {
                    this.enterEvent = undefined;
                });
            },
            tab(e) {
                this.tabEvent = e;
                this.$nextTick(() => {
                    this.tabEvent = undefined;
                });
            },
            esc(e) {
                if(this.columns === 0) {
                    this.close();
                }
                else {
                    this.escEvent = e;
                }

                this.$nextTick(() => {
                    this.escEvent = undefined;
                });
            },
            up(e) {
                this.upEvent = e;
                this.$nextTick(() => {
                    this.upEvent = undefined;
                });
            },
            down(e) {
                this.downEvent = e;
                this.$nextTick(() => {
                    this.downEvent = undefined;
                });
            },
            del(e) {
                this.delEvent = e;
                this.$nextTick(() => {
                    this.delEvent = undefined;
                });
            },
            back() {
                if(this.columns > 0) {
                    this.$store.dispatch("popColumn");
                    this.focusCmdInput();
                }
                else {
                    this.close();
                }
            },
            action() {
                this.$store.dispatch("resetCmd", "/");
            },
            navigation() {
                this.$store.dispatch("resetCmd", ".");
            },
            clear() {
                this.$store.dispatch("resetCmd");
            },
            close() {
                this.$store.commit("setCmdOpen", false);
                this.$store.dispatch("resetCmd");
            }
        }
    }
</script>
