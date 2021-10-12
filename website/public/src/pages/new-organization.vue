<template>
    <layout :active="'organization'">
        <template>
            <div class="section-reg">
                <div class="container">
                    <div class="reg-card center-xs">
                        <div class="progress">
                            <div class="icon icon-expand icon-rotate-90" v-on:click="stepBack" v-show="step > 0 && step < 3"></div>
                            <div v-bind:class="{'step': true, 'done': step >= 0}"></div>
                            <div v-bind:class="{'step': true, 'done': step >= 1}"></div>
                            <div v-bind:class="{'step': true, 'done': step >= 2}"></div>
                            <div v-bind:class="{'step': true, 'done': step >= 3}"></div>
                        </div>
                        <emvi-create-organization-name-url v-if="step === 0" v-on:next="nextStep"></emvi-create-organization-name-url>
                        <emvi-create-organization-user-lang v-if="step === 1" v-on:next="nextStep"></emvi-create-organization-user-lang>
                        <emvi-create-organization-confirmation v-if="step === 2" v-on:next="nextStep"></emvi-create-organization-confirmation>
                        <emvi-create-organization-completed v-if="step === 3" v-on:next="nextStep" :domain="domain"></emvi-create-organization-completed>
                        <div class="reg-footer">
                            <template v-if="step !== 3">
                                <router-link to="/organizations">{{$t("link_cancel")}}</router-link>
                            </template>
                            <template v-if="step === 3">
                                <router-link to="/organizations">{{$t("link_overview")}}</router-link>
                            </template>
                        </div>
                    </div>
                </div>
            </div>
        </template>
    </layout>
</template>

<script>
    import {
        layout,
        emviCreateOrganizationNameUrl,
        emviCreateOrganizationUserLang,
        emviCreateOrganizationConfirmation,
        emviCreateOrganizationCompleted
    } from "../components";

    export default {
        components: {
            layout,
            emviCreateOrganizationNameUrl,
            emviCreateOrganizationUserLang,
            emviCreateOrganizationConfirmation,
            emviCreateOrganizationCompleted
        },
        data() {
            return {
                step: 0,
                domain: ""
            };
        },
        methods: {
            nextStep(domain) {
                this.step++;
                this.domain = domain;
            },
            stepBack() {
                this.step--;
            }
        }
    }
</script>

<style scoped>
    .reg-card {
        text-align: left;
    }
</style>

<i18n>
    {
        "en": {
            "link_cancel": "Cancel",
            "link_overview": "Go to my organizations"
        },
        "de": {
            "link_cancel": "Abbrechen",
            "link_overview": "Zu meinen Organisationen"
        }
    }
</i18n>
