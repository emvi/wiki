<template>
	<modal v-on:action="uploadPicture" v-on:cancel="cancel" size="medium" :action-label="$t('button_upload')" :action-icon="actionIcon" :action-spinner="actionIcon === 'sync'" action-color="pink" :error="err">
		<template slot="title">
			{{$t("title")}}
		</template>
		<template>
			<fileinput v-model="file"></fileinput>
			<small class="form--element--hint">{{$t("hint_picture")}}</small>
			<btn icon="trash" color="red" v-on:click="deletePicture">{{$t("delete_picture")}}</btn>
		</template>
	</modal>
</template>

<script>
	import {AuthService, ErrorService} from "../service";
	import {btn, fileinput} from "../components";
	import modal from "./modal.vue";

	const MAX_SIZE_BYTES = 512000; // 512 Kb

	export default {
		components: {
			modal,
			btn,
			fileinput
		},
		data() {
			return {
				err: null,
				file: null,
				actionIcon: "img"
			}
		},
		methods: {
			uploadPicture() {
				let form = new FormData();
				form.append("file", this.file);

				if(this.file.size > MAX_SIZE_BYTES) {
					this.err = this.$t("err_max_size");
					return;
				}

				this.actionIcon = "sync";

				AuthService.uploadPicture(form)
				.then(() => {
					this.$store.dispatch("loadUser");
					this.$store.dispatch("showToast", {type: "green", message: this.$t("toast_uploaded")});
					this.$emit("action");
				})
				.catch(e => {
					this.actionIcon = "img";
					this.setError(e);
				});
			},
			deletePicture() {
				AuthService.deletePicture()
				.then(() => {
					this.$store.dispatch("loadUser");
					this.$store.dispatch("showToast", {type: "green", message: this.$t("toast_deleted")});
					this.$emit("action");
				})
				.catch(e => {
					this.setError(e);
				});
			},
			cancel() {
				this.$emit("close");
			}
		}
	}
</script>

<i18n>
	{
		"en": {
			"title": "Change Profile Picture",
			"delete_picture": "Delete Picture",
			"button_upload": "Upload Picture",
			"hint_picture": "The maximum size for the picture is 512 KB",
			"toast_uploaded": "Picture updated.",
			"toast_deleted": "Picture deleted.",
			"err_max_size": "The file exceeds the maximum size."
		},
		"de": {
			"title": "Profilbild ändern",
			"delete_picture": "Bild löschen",
			"button_upload": "Bild hochladen",
			"hint_picture": "Die maximale Größe für das Bild beträgt 512 KB",
			"toast_uploaded": "Bild aktualisiert.",
			"toast_deleted": "Bild gelöscht.",
			"err_max_size": "Die Datei überschreitet die maximale Größe."
		}
	}
</i18n>
