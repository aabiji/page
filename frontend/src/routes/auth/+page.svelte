<script>
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

    function validateAuthInfo() {
        if (authInfo.email == "" || authInfo.password == "" ||
            (authInfo.confirm == "" && !isLogin)) {
            authError = "Please fill out all form fields.";
        } else if (!authInfo.email.match(/^\S+@\S+\.\S+$/)) {
            authError = "Please enter a valide email address.";
        } else if (!isLogin && authInfo.confirm != authInfo.password) {
            authError = "Password and repeated password must match";
        } else {
            authError = "";
        }
    }

    async function authenticate() {
        validateAuthInfo();
        if (authError != "") return;
        let url = `http://localhost:8080/user/${isLogin ? "auth" : "create"}`;
        let unhashedPassword = authInfo.password;
        authInfo.password = await utils.hashSHA256(authInfo.password);
        utils.callApi(url, "POST", authInfo).then((response) => {
            authInfo.password = unhashedPassword;
            if (utils.serverError in response) {
                authError = response[utils.serverError];
                return;
            }
            goto("/read");
        });
    }

    onMount(() => {
        utils.redirectIfNotAuth();
    });
</script>

<div class="container">
    <div>
        <h2 class="logo"> Page </h2>
        <h2> | {authState} </h2><br>
        <p class="error-message"> {authError} </p>

        <input bind:value={authInfo.email} type="email" placeholder="Email"><br>
        <input bind:value={authInfo.password} type="password" placeholder="Password"><br>
        {#if !isLogin}
            <input bind:value={authInfo.confirm} type="password" placeholder="Repeat password"><br>
        {/if}
        <button on:click={authenticate}> {authState} </button><br>

        {#if !isLogin}
            <a on:click={toggleState}> Already have an account? </a><br>
        {:else}
            <a> Forgot password? </a><br>
            <a on:click={toggleState}> Don't have an account? </a><br>
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
    a {
        float: right;
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
    button {
        width: 325px;
        color: white;
        font-size: 18px;
        padding: 10px 10px;
        margin-bottom: 15px;
        background-color: var(--accent-color);
    }
    button:hover {
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
</style>
