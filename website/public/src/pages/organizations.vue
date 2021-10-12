<template>
	<layout :active="'organization'">
		<template>
			<div style="margin-top: 48px;">
				<transition-group name="fade-in">
					<invitationcard v-for="inv in invitations"
						:key="inv.id"
						:invitation="inv"
						v-on:delete="deleteInvitation"></invitationcard>
				</transition-group>
				<div class="spacer-32" v-if="invitations.length && organizations.length && loaded"></div>
				<transition-group name="fade-in" v-show="loaded">
					<organizationcard v-for="org in organizations"
						:key="org.id"
						:orga="org"></organizationcard>
				</transition-group>
				<template v-if="!invitations.length && !organizations.length && loaded">
					<div class="modal--warning">{{$t("hint_invite")}}</div>
					<div class="spacer-16"></div>
				</template>
				<div class="large-card--placeholder" v-show="!loaded">
					<div class="large-card--pic"></div>
					<div class="large-card--content">
						<div class="large-card--content--title"></div>
						<div class="large-card--content--info"></div>
					</div>
				</div>
			</div>
			<router-link to="/new" class="button--large" v-if="freeOrganization">
				<i class="icon large-icon--grey icon-add"></i>
				<div class="button--large--text">
					{{$t("button_new_organization")}}
				</div>
			</router-link>
		</template>
	</layout>
</template>

<script>
	import {OrganizationService, InvitationService} from "../service";
	import {title} from "../mixins";
	import {organizationcard, layout, invitationcard} from "../components";

	export default {
		mixins: [title],
		components: {
			organizationcard,
			invitationcard,
			layout
		},
		data() {
			return {
				organizations: [],
				invitations: [],
				loaded: false
			}
		},
        computed: {
		    freeOrganization() {
		        let hasOneFree = true;

		        for(let i = 0; i < this.organizations.length; i++) {
		            if(!this.organizations[i].expert) {
		                hasOneFree = false;
                        break;
                    }
                }

		        return hasOneFree;
            }
        },
        beforeMount() {
		    let state = this.$route.query.state;

		    if(state) {
                this.redirect(state);
            }
        },
        mounted() {
			this.loadOrganizations();
			this.loadInvitations();
		},
		methods: {
			loadOrganizations() {
				OrganizationService.getOrganizations()
				.then(organizations => {
					this.organizations = organizations;
					this.loaded = true;
				});
			},
            loadInvitations() {
				InvitationService.getInvitations()
				.then(invitations => {
					this.invitations = invitations;
				});
			},
            deleteInvitation(id) {
				InvitationService.deleteInvitation(id)
				.then(() => {
					this.loadInvitations();
				});
			},
            redirect(name) {
                let url = EMVI_WIKI_FRONTEND_HOST.split("/");
                window.location = `${url[0]}//${name}.${url[2]}`;
            }
        }
	}
</script>

<i18n>
	{
		"en": {
			"pagetitle": "Organizations",
			"hint_invite": "You are currently not a member of an organization. Ask an administrator to invite you to an existing organization or create a new one below",
			"button_new_organization": "Create New Organization"
		},
		"de": {
			"pagetitle": "Organisationen",
			"hint_invite": "Du bist momentan kein Mitglied einer Organisation. Bitte einen Administrator dich in eine existierende Organisation einzuladen oder lege eine Neue an.",
			"button_new_organization": "Neue Organisation anlegen"
		}
	}
</i18n>
