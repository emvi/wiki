<template>
	<transition appear name="slide">
		<div class="modal">
			<div v-bind:class="modalClasses" ref="modalBody">
				<div class="modal--top no-select">
					<slot name="title"></slot>
				</div>
				<transition name="fade">
					<div class="modal--error" v-show="error">
						{{error}}
					</div>
				</transition>
				<div class="modal--content">
					<slot></slot>
				</div>
				<div class="modal--bottom no-select">
					<btn type="solid" :color="cancelColor" v-on:click="$emit('cancel')">{{cancelButton}}</btn>
					<btn type="solid" :icon="actionIcon" :color="actionColor" :disabled="actionDisabled" :spinner="actionSpinner" v-on:click="$emit('action')" v-if="actionLabel">{{actionButton}}</btn>
				</div>
			</div>
		</div>
	</transition>
</template>

<script>
	import {btn} from "../components";

	export default {
		components: {btn},
		props: {
			size: {default: "small"}, // small, medium, large
			title: null,
			actionLabel: null,
			actionIcon: null,
			actionColor: {default: "blue"},
			actionDisabled: {default: false},
			actionSpinner: {default: false},
			cancelLabel: null,
			cancelColor: {default: "grey"},
			close: {default: true},
			error: {default: false}
		},
		watch: {
			actionLabel(value) {
				this.actionButton = value;
			}
		},
		data() {
			return {
				modalClasses: {},
				actionButton: "",
				cancelButton: ""
			};
		},
		created() {
			this.modalClasses["modal--card modal--card--"+this.size] = true;

			if(!this.actionLabel) {
				this.actionButton = this.$t("button_save");
			}
			else {
				this.actionButton = this.actionLabel;
			}

			if(!this.cancelLabel) {
				this.cancelButton = this.$t("button_cancel");
			}
			else {
				this.cancelButton = this.cancelLabel;
			}

			document.addEventListener("mousedown", this.documentClick);
			document.addEventListener("keyup", this.keyup);
		},
		beforeDestroy() {
			document.removeEventListener("mousedown", this.documentClick);
			document.removeEventListener("keyup", this.keyup);
		},
		methods: {
			documentClick(e){
				let modal = this.$refs.modalBody;
				let target = e.target;

				if(modal && modal !== target && !modal.contains(target)){
					this.$emit("cancel");
				}
			},
			keyup(e) {
				if(e.keyCode === 27) {
					this.$emit("cancel");
				}
			}
		}
	}
</script>
