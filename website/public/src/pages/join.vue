<template>
	<layout>
		<template>
			<div class="max-width-m" style="margin: 48px auto 0 auto">
				<div class="row">
					<h3>{{$t("title")}}</h3>
					<div class="spacer-16"></div>
				</div>
				<div v-if="loggedIn">
					<div v-if="!notfound">
						<div class="spacer-16"></div>
						<h4>{{$t("organization_name")}}: {{organization.name}}</h4>
						<div class="spacer-16"></div>
						<p>{{$t("hint_join")}}</p>
						<div class="spacer-16"></div>
						<form v-on:submit.prevent="join">
							<field v-model="username" :type="'text'" :label="$t('label_username')" :error="validation['username']" color="pink"></field>
						</form>
						<btn type="solid" icon="arrow" v-on:click="join">{{$t("button_join_organization")}}</btn>
					</div>
					<div v-if="notfound">
						<p>{{$t("invitation_not_found")}}</p>
						<p>
							{{$t("invitation_not_found_mail")}}
							<strong>{{email}}</strong>
						</p>
					</div>
				</div>
				<div v-if="!loggedIn">
					<p>{{$t("hint_login")}}</p>
					<div class="spacer-16"></div>
					<div class="row">
						<btn type="solid" v-on:click="register">{{$t("button_register")}}</btn>
						<btn type="solid" color="grey" v-on:click="login">{{$t("button_login")}}</btn>
					</div>
				</div>
			</div>
		</template>
	</layout>
</template>

<script>
	import {MemberService} from "../service";
	import {title} from "../mixins";
	import {layout, mainmenu, field, btn} from "../components";

	export default {
		mixins: [title],
		components: {layout, mainmenu, field, btn},
		data() {
			return {
				notfound: false,
				code: "",
				invitationCode: "",
				organization: {},
				username: ""
			}
		},
		computed: {
			loggedIn() {
				return this.$store.getters.loggedIn;
			},
			email() {
				return this.$store.state.user.user.email;
			}
		},
		mounted() {
			this.code = this.$route.query.code;
			this.invitationCode = this.$route.params.code;

			if(this.$store.getters.loggedIn){
				this.deleteInvitationCache();
				this.loadInvitation();
			}
			else{
				this.setInvitationCache();
			}
		},
		methods: {
			deleteInvitationCache() {
				window.localStorage.removeItem("join_code");
			},
			setInvitationCache() {
				window.localStorage.setItem("join_code", this.code);
			},
			loadInvitation() {
				if(this.invitationCode) {
					MemberService.getInvitationOrganization(this.invitationCode)
						.then(organization => {
							this.organization = organization;
						})
						.catch(() => {
							this.notfound = true;
						});
				}
				else {
					MemberService.getInvitation(this.code)
						.then(organization => {
							this.organization = organization;
						})
						.catch(() => {
							this.notfound = true;
						});
				}
			},
			login() {
				this.$store.dispatch("login");
			},
			register() {
				let email = this.$route.query.email;

				if(email) {
					this.$router.push(`/registration?email=${email}`);
				}
				else {
					this.$router.push("/registration");
				}
			},
			join() {
				MemberService.joinOrganization(this.username, this.code, this.invitationCode)
					.then(() => {
						this.$router.push("/organizations");
					})
					.catch(e => {
						this.setError(e);
					});
			}
		}
	}
</script>

<i18n>
	{
		"en": {
			"title": "Joining",
			"organization_name": "Organization",
			"label_username": "Username",
			"button_join_organization": "Join the Organization",
			"invitation_not_found": "Invitation not found. Have you been invited for a different account?",
			"invitation_not_found_mail": "Your are logged in as:",
			"hint_login": "You must be logged in or create a new account to join an organization.",
			"hint_join": "Choose a short and meaningful username. It cannot be changed!",
			"button_login": "Login",
			"button_register": "Create new account"
		},
		"de": {
			"title": "Beitritt",
			"organization_name": "Organisation",
			"label_username": "Nutzername",
			"button_join_organization": "Organisation beitreten",
			"invitation_not_found": "Die Einladung konnte nicht gefunden werden. Wurdest du vielleicht für einen anderen Account eingeladen?",
			"invitation_not_found_mail": "Du bist eingeloggt als:",
			"hint_login": "Bevor Du der Organisation beitreten kannst, musst Du dich anmelden oder registrieren.",
			"hint_join": "Wähle einen kurzen und prägnanten Nutzernamen. Er kann nicht mehr geändert werden!",
			"button_login": "Anmelden",
			"button_register": "Registrieren"
		}
	}
</i18n>
