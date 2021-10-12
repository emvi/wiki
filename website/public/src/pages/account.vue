<template>
	<layout :active="'settings'">
		<uploadpicture v-if="showUpload"
			v-on:action="showUpload = false"
			v-on:close="showUpload = false"></uploadpicture>

		<template>
			<div class="settings--account">
				<div class="settings--account--pic">
					<img v-bind:src="picture" alt="" v-if="picture" />
					<i class="icon icon-img" v-if="!picture"></i>
					<div class="settings--account--pic--edit" v-on:click.stop="showUpload = true">
						<i class="icon icon-edit"></i>
					</div>
				</div>
				<div class="settings--account--name">
					{{user.firstname}} {{user.lastname}}
				</div>
			</div>
			<card collapsible="true" color="pink" :maxheight="false" ref="namecard" show-info="true" show-title="true">
				<template slot="icon">
					<i class="icon large-icon--grey icon-user card--content--icon"></i>
				</template>
				<template slot="info">
					{{$t("title_name")}}
				</template>
				<template slot="title">
					{{user.firstname}} {{user.lastname}}
				</template>
				<template slot="children">
					<div class="card--content--settings">
						<div class="form--element--full">
							<form v-on:submit.prevent="changeNameData">
								<field v-model="firstname" :type="'text'" :label="$t('label_firstname')" color="pink" :error="validation['firstname']"></field>
								<field v-model="lastname" :type="'text'" :label="$t('label_lastname')" color="pink" :error="validation['lastname']"></field>
								<button type="submit"></button>
							</form>
						</div>
					</div>
					<div class="card--content--buttons">
						<btn type="solid" color="grey" v-on:click="reset">{{$t('button_discard')}}</btn>
						<btn type="solid" icon="save" color="pink" v-on:click="changeNameData">{{$t('button_save')}}</btn>
					</div>
				</template>
			</card>
			<card :collapsible="!isSSOUser" :cursor-default="isSSOUser" color="pink" :maxheight="false" ref="emailcard" show-info="true" show-title="true">
				<template slot="icon">
					<i class="icon large-icon--grey icon-mail card--content--icon"></i>
				</template>
				<template slot="info">
					{{$t("title_mail")}}
				</template>
				<template slot="title">
					{{user.email}}
				</template>
				<template slot="children">
					<div class="card--content--settings">
						<div class="form--element--full">
							<form v-on:submit.prevent="changeMail">
								<field v-model="email" type="email" :label="$t('label_email')" :error="validation['email']" color="pink"></field>
								<button type="submit"></button>
							</form>
						</div>
					</div>
					<div class="card--content--buttons">
						<btn type="solid" color="grey" v-on:click="reset">{{$t('button_discard')}}</btn>
						<btn type="solid" color="pink" icon="save" v-on:click="changeMail">{{$t('button_save')}}</btn>
					</div>
				</template>
			</card>
			<card collapsible="true" color="pink" :maxheight="false" ref="passwordcard" v-if="!isSSOUser" show-info="true" show-title="true">
				<template slot="icon">
					<i class="icon large-icon--grey icon-lock card--content--icon"></i>
				</template>
				<template slot="info">
					{{$t("title_password")}}
				</template>
				<template slot="title">
					&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;
				</template>
				<template slot="children">
					<div class="card--content--settings">
						<div class="form--element--full">
							<form v-on:submit.prevent="changePassword">
								<field v-model="old_password" :type="'password'" :label="$t('label_old_pwd')" :error="validation['old_password']" color="pink"></field>
								<field v-model="new_password" :type="'password'" :label="$t('label_new_pwd')" :error="validation['password']" color="pink"></field>
								<field v-model="new_password_repeat" :type="'password'" :label="$t('label_new_pwd_repeat')" color="pink"></field>
								<button type="submit"></button>
							</form>
						</div>
					</div>
					<div class="card--content--buttons">
						<btn type="solid" icon="save" color="pink" v-on:click="changePassword">{{$t('button_save')}}</btn>
					</div>
				</template>
			</card>
			<card collapsible="true" color="pink" :maxheight="false" ref="langcard" show-info="true" show-title="true">
				<template slot="icon">
					<i class="icon large-icon--grey icon-language card--content--icon"></i>
				</template>
				<template slot="info">
					{{$t("title_language")}}
				</template>
				<template slot="title" v-if="!language">
					{{$t('label_no_language')}}
				</template>
				<template slot="title" v-if="language">
					{{userLanguage}} ({{user.language}})
				</template>
				<template slot="children">
					<div class="card--content--settings">
						<div class="form--element--full">
							<selection v-model="language" :label="$t('label_language')" :hint="$t('hint_language')"> 
								<option v-for="l in langs" :value="l.code" :key="l.code">{{l.nativeName}}</option>
							</selection>
						</div>
					</div>
					<div class="card--content--buttons">
						<btn type="solid" icon="save" color="pink" v-on:click="changeLanguage">{{$t('button_save')}}</btn>
					</div>
				</template>
			</card>
			<card collapsible="true" color="pink" :maxheight="false" ref="colormodecard" show-info="true" show-title="true">
				<template slot="icon">
					<i class="icon large-icon--grey icon-bgcolor card--content--icon"></i>
				</template>
				<template slot="info">
					{{$t("title_color_scheme")}}
				</template>
				<template slot="title">
					<span v-show="parseInt(colormode) === 0">{{$t("label_color_default")}}</span>
					<span v-show="parseInt(colormode) === 1">{{$t("label_color_protanopia_deuteranopia")}}</span>
					<span v-show="parseInt(colormode) === 2">{{$t("label_color_tritanopia")}}</span>
				</template>
				<template slot="children">
					<div class="card--content--settings">
						<div class="form--element--radio">
							<label class="form--element--radio--label">
								<input name="colormode" class="form--element--radio--field" type="radio" value="0" v-model="colormode" />
								<i class="form--element--radio--button"></i>
								<span>{{$t("label_color_default")}}</span>
							</label>
							<label class="form--element--radio--label">
								<input name="colormode" class="form--element--radio--field" type="radio" value="1" v-model="colormode" />
								<i class="form--element--radio--button"></i>
								<span>{{$t("label_color_protanopia_deuteranopia")}}</span>
							</label>
							<label class="form--element--radio--label">
								<input name="colormode" class="form--element--radio--field" type="radio" value="2" v-model="colormode" />
								<i class="form--element--radio--button"></i>
								<span>{{$t("label_color_tritanopia")}}</span>
							</label>
							<div class="no-select">
								<div class="colors background--blue"></div>
								<div class="colors background--purple"></div>
								<div class="colors background--pink"></div>
								<div class="colors background--red"></div>
								<div class="colors background--orange"></div>
								<div class="colors background--yellow"></div>
								<div class="colors background--green"></div>
							</div>
						</div>
					</div>
				</template>
			</card>
			<card>
				<template slot="icon">
					<i class="icon large-icon--grey icon-info card--content--icon cursor-default"></i>
				</template>
				<template>
					<label class="form--element--checkbox--label">
						<input class="form--element--checkbox--field" type="checkbox" v-model="marketing" v-on:click="changeMarketing" />
						<i class="form--element--checkbox--button"></i>
						<span>{{$t("label_newsletter")}}</span>
					</label>
				</template>
			</card>
		</template>
	</layout>
</template>

<script>
	import ISO6391 from "iso-639-1";
	import {mapGetters} from "vuex";
	import {AuthService} from "../service";
	import {title} from "../mixins";
	import {uploadpicture} from "../modal";
	import {mainmenu, field, layout, selection, card, btn} from "../components";
	import {getSupportedLocale} from "../util/locale";

	export default {
		mixins: [title],
		components: {
			mainmenu,
			uploadpicture,
			field,
			layout,
			selection,
			card,
			btn
		},
		computed: {
		    ...mapGetters(["user"]),
			langs() {
				return ISO6391.getLanguages(ISO6391.getAllCodes());
			},
			userLanguage() {
				return ISO6391.getNativeName(this.user.language);
			},
			picture() {
				return this.user.picture;
			},
			isSSOUser() {
				return this.user.is_sso_user;
			}
		},
        watch: {
            colormode(value) {
                if(value !== this.userSettings.colormode) {
                    AuthService.saveColorMode(value)
                        .then(() => {
                            this.validation = {};
                            this.userSettings.colormode = parseInt(this.colormode);
                            this.$store.dispatch("loadUser");
                            this.$store.dispatch("showToast", {type: "green", message: this.$t("toast_saved")});
                        });
                }
            }
        },
		data() {
			return {
				showUpload: false,
				userSettings: {},
				language: "",
				marketing: false,
				firstname: "",
				lastname: "",
				email: "",
				old_password: "",
				new_password: "",
				new_password_repeat: "",
				colormode: 0
			}
		},
		mounted() {
			this.userSettings = {
				firstname: this.user.firstname,
				lastname: this.user.lastname,
				language: this.user.language,
				email: this.user.email,
				marketing: this.user.accept_marketing,
				colormode: this.user.color_mode
			};
			this.reset();
		},
		methods: {
			changeNameData() {
				AuthService.savePersonalData(this.firstname, this.lastname, this.user.language, this.user.marketing)
				.then(() => {
					this.validation = {};
					this.userSettings.firstname = this.firstname;
					this.userSettings.lastname = this.lastname;
					this.$refs.namecard.collapse();
					this.$store.dispatch("loadUser");
					this.$store.dispatch("showToast", {type: "green", message: this.$t("toast_saved")});
				})
				.catch(e => {
					this.setError(e);
				});
			},
			changeMail() {
				AuthService.saveMail(this.email)
				.then(() => {
					this.validation = {};
					this.$store.dispatch("showToast", {type: "green", message: this.$t("toast_mail")});
				})
				.catch(e => {
					this.setError(e);
				});
			},
			changePassword() {
				AuthService.savePassword(this.old_password, this.new_password, this.new_password_repeat)
				.then(() => {
					this.validation = {};
					this.old_password = "";
					this.new_password = "";
					this.new_password_repeat = "";
					this.$refs.passwordcard.collapse();
					this.$store.dispatch("loadUser");
					this.$store.dispatch("showToast", {type: "green", message: this.$t("toast_saved")});
				})
				.catch(e => {
					this.setError(e);
				});
			},
			changeLanguage() {
				AuthService.savePersonalData(this.user.firstname, this.user.lastname, this.language, this.user.marketing)
				.then(() => {
					this.validation = {};
					this.$store.dispatch("loadUser");
					this.userSettings.language = this.language;
					this.$root.$i18n.locale = getSupportedLocale(this.language);
					this.$refs.langcard.collapse();
					this.$store.dispatch("showToast", {type: "green", message: this.$t("toast_saved")});
				})
				.catch(e => {
					this.setError(e);
				});
			},
			changeMarketing() {
				AuthService.savePersonalData(this.user.firstname, this.user.lastname, this.user.language, !this.marketing)
				.then(() => {
					this.validation = {};
					this.$store.dispatch("loadUser");
					this.userSettings.marketing = this.marketing;
					this.$store.dispatch("showToast", {type: "green", message: this.$t("toast_saved")});
				})
				.catch(e => {
					this.setError(e);
				});
			},
			reset() {
				this.validation = {};
				this.language = this.userSettings.language;
				this.marketing = this.userSettings.marketing;
				this.email = this.userSettings.email;
				this.firstname = this.userSettings.firstname;
				this.lastname = this.userSettings.lastname;
				this.colormode = this.userSettings.colormode;
			}
		}
	}
</script>

<i18n>
	{
		"en": {
			"pagetitle": "Account",
			"title_name": "Name",
			"title_mail": "Email",
			"title_password": "Password",
			"title_language": "Prefered Language",
			"title_color_scheme": "Color Scheme",
			"label_firstname": "Firstname",
			"label_lastname": "Lastname",
			"label_email": "New Email",
			"label_old_pwd": "Old Password",
			"label_new_pwd": "New Password",
			"label_new_pwd_repeat": "Repeat New Password",
			"label_color_default": "Default",
			"label_color_protanopia_deuteranopia": "Protanopia & Deuteranopia",
			"label_color_tritanopia": "Tritanopia",
			"label_language": "Prefered Language",
			"label_no_language": "Unset",
			"hint_language": "Prefered Language when opening Articles. Furthermore sets the language of the user interface, currently English or German.",
			"label_newsletter": "I would like to receive regularly updates about Emvi via email.",
			"toast_mail": "A email was send to you to confirm your changes."
		},
		"de": {
			"pagetitle": "Account",
			"title_name": "Name",
			"title_mail": "E-Mail-Adresse",
			"title_password": "Passwort",
			"title_language": "Bevorzugte Sprache",
			"title_color_scheme": "Farbschema",
			"label_firstname": "Vorname",
			"label_lastname": "Nachname",
			"label_email": "Neue E-Mail-Adresse",
			"label_old_pwd": "Altes Passwort",
			"label_new_pwd": "Neues Passwort",
			"label_new_pwd_repeat": "Neues Passwort wiederholen",
			"label_language": "Bevorzugte Sprache",
			"label_no_language": "Nicht gesetzt",
			"label_color_default": "Standard",
			"label_color_protanopia_deuteranopia": "Rotblindheit & Gründblindheit",
			"label_color_tritanopia": "Blaublindheit",
			"hint_language": "Bevorzugte Sprache beim Öffnen von Artikeln. Legt außerdem die Sprache der Oberfläche fest, momentan Deutsch oder Englisch.",
			"label_newsletter": "Ich möchte regelmäßig Informationen zu Emvi per E-Mail erhalten.",
			"toast_mail": "Dir wurde eine E-Mail zur Bestätigung der Änderung gesendet."
		}
	}
</i18n>
