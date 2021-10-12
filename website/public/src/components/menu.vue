<template>
	<div class="top-bar no-select">
		<div class="main-grid">
			<div v-bind:class="{'main-grid--side': $mq > 1024, 'top-bar--left': true}">
				<a href="/" class="top-bar--wordmark">
					<img v-bind:src="'./../static/img/wordmark-black.svg'" alt="Emvi" />
				</a>
			</div>
			<div class="main-grid--center top-bar--hub-menu" v-if="user && $mq > 768">
				<router-link to="/account">
					<btn icon="user" color="pink" :active="active === 'settings'">{{$t("button_account")}}</btn>
				</router-link>
				<router-link to="/organizations">
					<btn icon="organization" :active="active === 'organization'">{{$t("button_organizations")}}</btn>
				</router-link>
			</div>
			<div class="main-grid--center top-bar--hub-menu" v-if="!user && $mq > 768">
				<ul class="col-sm-6 center-sm header--center"></ul>
			</div>
			<div class="main-grid--side top-bar--hub-side-menu" v-if="user && $mq > 768">
				<btn icon="logout" v-on:click="logout">{{$t("link_logout")}}</btn>
			</div>
			<div class="main-grid--side top-bar--hub-side-menu" v-if="!user && $mq > 768">
				<div class="col-sm-3 col-xs-8 end-xs header--right">
					<a href="#" class="header--secondary" v-on:click.prevent="login">
						{{$t("link_login")}}
					</a>
					<router-link to="/registration" class="header--primary">
						{{$t("link_register")}}
					</router-link>
				</div>
			</div>
			<div class="top-bar--menu top-bar--menu" v-if="user && $mq <= 768">
				<router-link to="/account">
					<i v-bind:class="{'top-bar--menu--icon icon icon-user': true, 'color--pink-100': active === 'settings'}"></i>
				</router-link>
				<router-link to="/organizations">
					<i v-bind:class="{'top-bar--menu--icon icon icon-organization': true, 'color--blue-100': active === 'organization'}"></i>
				</router-link>
				<i class="top-bar--menu--icon icon icon-logout" v-on:click="logout"></i>
			</div>
		</div>
	</div>
</template>

<script>
	import btn from "./html/btn.vue";

	export default {
		components: {btn},
		props: ["active"],
		data() {
			return {
				usermenu_visible: false,
				notifications_visible: false
			}
		},
		computed: {
			user() {
				return this.$store.state.user.user;
			}
		},
		methods: {
			login() {
				this.$store.dispatch("login");
			},
			logout() {
				this.$store.dispatch("logout");
			}
		}
	}
</script>

<i18n>
	{
		"en": {
			"button_account": "Account",
			"button_organizations": "Organizations",
			"link_login": "Login",
			"link_logout": "Logout",
			"link_register": "Create Account"
		},
		"de": {
			"button_account": "Konto",
			"button_organizations": "Organisationen",
			"link_login": "Anmelden",
			"link_logout": "Abmelden",
			"link_register": "Registrieren"
		}
	}
</i18n>
