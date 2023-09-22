<script>
    // TODO: login with different service providers
    let isLogin = true;
    const toggleState = () => isLogin = !isLogin;
</script>

{#if isLogin}
<div class="container">
    <div>
        <h2 class="logo"> Page </h2>
        <h2> | Login </h2><br><br>
        <input type="email" placeholder="Email"><br>
        <input type="password" placeholder="Password"><br>
        <button> Login </button><br>
        <a> Forgot password? </a><br>
        <a on:click={toggleState}> Don't have an account? </a><br>
    </div>
</div>
{:else}
<div class="container">
    <div>
        <h2 class="logo"> Page </h2>
        <h2> | Create account </h2><br><br>
        <input type="email" placeholder="Email"><br>
        <input type="password" placeholder="Password"><br>
        <input type="password" placeholder="Repeat password"><br>
        <button> Signup </button><br>
        <a on:click={toggleState}> Already have an account? </a><br>
    </div>
</div>
{/if}

<style>
    .logo {
        width: 300px;
        font-size: 30px;
        margin-right: 5px;
    }
    h2 {
        display: inline;
        margin-bottom: 15px;
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
        background-color: #1b1c1c;
        transform: translate(-50%, -50%);
        box-shadow: rgba(0, 0, 0, 0.35) 0px 5px 15px;
    }
</style>
