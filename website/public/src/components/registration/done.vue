<template>
    <div class="reg-form">
        <h3 class="reg-headline">
            {{$t("title")}}
        </h3>
        <div class="reg-icon">
		    <div class="icon icon-tick"></div>
		</div>
        <p>
            {{$t("text")}}
        </p>
        <div class="row">
            <button class="col-sm-12 reg-button--primary" v-on:click="joinOrga">{{$t("button_join")}}</button>
            <button class="col-sm-12 reg-button--secondary" v-on:click="newOrga">{{$t("button_new")}}</button>
        </div>
    </div>
</template>

<script>
    export default {
        methods: {
            newOrga() {
                this.$router.push("/new");
            },
            joinOrga() {
                this.$router.push("/organizations");
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "title": "Registration completed!",
            "text": "Welcome to Emvi! You can now create your own organization or join an existing one.",
            "button_new": "Create Organization",
            "button_join": "Join Organization"
        },
        "de": {
            "title": "Registrierung abgeschlossen!",
            "text": "Willkommen bei Emvi! Du kannst nun deine eigene Organisation anlegen oder einer Bestehenden beitreten.",
            "button_new": "Organisation anlegen",
            "button_join": "Organisation beitreten"
        }
    }
</i18n>
