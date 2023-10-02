<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import * as utils from "$lib/utils";

    let isLogin = true;
    let authError = "";
    let authState = "Login"
    let authInfo = {email: "", password: "", confirm: ""};
    const toggleState = () => {
        authError = "";
        isLogin = !isLogin;
        authInfo = {email: "", password: "", confirm: ""};
        authState = isLogin ? "Login" : "Create account";
    }

    function isValidEmail(): boolean {
        let regex = /^\S+@\S+\.\S+$/;
        return regex.test(authInfo.email);
    }
 
    function isIncompleteForm(): boolean {
        let confirmEmpty = !isLogin && authInfo.confirm == "";
        return authInfo.email == "" || authInfo.password == "" || confirmEmpty;
    }

    // To be "safe", a password must be at least 10 characters long
    // and must contain at least 1 special character.
    function isSafePassword(): boolean {
        let longEnough = authInfo.password.length >= 10;
        let specialChars = ["!", "~", "@", "#", "$", "%", "&", "*", "^", "?"];
        let hasSpecialChar = specialChars.some((c) => authInfo.password.includes(c));
        return longEnough && hasSpecialChar;
    }

    function validateAuthInfo() {
        if (isIncompleteForm()) {
            authError = "Please fill out all form fields.";
        } else if (!isValidEmail()) {
            authError = "Please enter a valid email address.";
        } else if (!isSafePassword()) {
            authError = "Password must be at least 8 characters long and must contain 1 special character.";
        } else if (!isLogin && authInfo.password != authInfo.confirm) {
            authError = "Password and repeated password must match.";
        } else {
            authError = "";
        }
    }

    async function hashSHA256(data: string): Promise<string> {
        let encoded = new TextEncoder().encode(data);
        let buffer = await window.crypto.subtle.digest("SHA-256", encoded);
        let hash = Array.from(new Uint8Array(buffer));
        return hash.map(byte => byte.toString(16).padStart(2, "0")).join("");
    }

    async function authenticate() {
        validateAuthInfo();
        if (authError != "") return;
        let url = `${utils.backendOrigin}/user/${isLogin ? "login" : "create"}`;
        let unhashedPassword = authInfo.password;
        authInfo.password = await hashSHA256(authInfo.password);
        utils.callApi(url, "POST", authInfo).then((response) => {
            authInfo.password = unhashedPassword;
            if (utils.serverError in response) {
                authError = response[utils.serverError];
                return;
            }
            goto("/");
        });
    }
 
    onMount(() => {
        utils.redirectIfNotAuth();
        // Submit form with enter key
        document.onkeyup = (event) => {
            if (event.key != "Enter") return;
            let submit = document.getElementsByClassName("button")[0] as HTMLElement;
            submit.click();
        }
    });
</script>

<div class="container">
    <div class="inner">
        <h2 class="logo"> Page </h2>
        <h2> | {authState} </h2><br>
        <p class="error-message"> {authError} </p>

        <input bind:value={authInfo.email} type="email" placeholder="Email"><br>
        <input bind:value={authInfo.password} type="password" placeholder="Password"><br>
        {#if !isLogin}
            <input bind:value={authInfo.confirm} type="password" placeholder="Repeat password"><br>
        {/if}
        <button class="button" on:click={authenticate}> {authState} </button><br>

        {#if !isLogin}
            <button class="option" on:click={toggleState}> Already have an account? </button>
        {:else}
            <button class="option"> Forgot password? </button><br>
            <button class="option" on:click={toggleState}> Don't have an account? </button>
        {/if}
    </div>
</div>

<style>
    .error-message {
        color: #ed3f2f;
    }
    .logo {
        width: 300px;
        font-size: 30px;
        margin-right: 5px;
    }
    h2 {
        display: inline;
    }
    .option {
        float: right;
        font-size: 14px;
        color: var(--accent-color);
        background-color: var(--background-accent);
    }
    input {
        width: 300px;
        color: white;
        font-size: 18px;
        padding: 10px 10px;
        margin-bottom: 15px;
        border: #535454 1px solid;
        background-color: rgba(0,0,0,0);
    }
    input:hover, input:focus {
        outline: none;
        border: var(--accent-color) 1px solid;
    }
    input:-webkit-autofill, input:-webkit-autofill:hover, input:-webkit-autofill:focus {
        font-size: 18px;
        -webkit-text-fill-color: white;
        -webkit-box-shadow: 0 0 0px 40rem #1b1c1c inset;
    }
    .button {
        width: 325px;
        color: white;
        font-size: 18px;
        padding: 10px 10px;
        margin-bottom: 15px;
        background-color: var(--accent-color);
    }
    .button:hover {
        background-color: var(--accent-color-darken);
    }
    .container {
        top: 50%;
        left: 50%;
        width: 600px;
        height: 600px;
        display: flex;
        margin: 0 auto;
        position: absolute;
        text-align: center;
        align-items: center;
        margin-bottom: 10px;
        justify-content: center;
        background-color: var(--background-accent);
        transform: translate(-50%, -50%);
        box-shadow: rgba(0, 0, 0, 0.35) 0px 5px 15px;
    }
    .inner {
        max-width: 350px;
    }
</style>