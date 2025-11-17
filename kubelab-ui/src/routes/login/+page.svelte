<script lang="ts">
  import { goto } from "$app/navigation";
  import ToggleConfetti from "$lib/components/base/ToggleConfetti.svelte";
  import { login, loginWithOAuth2 } from "$lib/pocketbase";
  import { alertOnFailure } from "$lib/pocketbase/ui";
  import darkTheme from "$lib/stores/theme";
  import type { PageData } from "./$types";

  // @ts-ignore
  import { Confetti } from "svelte-confetti";
  import toast from "svelte-french-toast";

  export let data: PageData;
  const { allowEmailAuth, allowOAuth2Auth } = data;

  const DEFAULTS = {
    email: "",
    password: ""
  };
  let user = { ...DEFAULTS };
  let loading = false;
  let oauthLoading = false;

  async function submit() {
    loading = true;
    await alertOnFailure(async function () {
      await login(user.email, user.password);
      toast.success("Logged in successfully!");
      goto("/app");
    }).finally(() => {
      loading = false;
    });
  }

  async function handleOAuthLogin() {
    oauthLoading = true;
    try {
      await loginWithOAuth2("github");
      toast.success("Logged in with GitHub!");
      goto("/app");
    } catch (error) {
      console.error("OAuth login failed:", error);
      toast.error("Failed to login with GitHub");
    } finally {
      oauthLoading = false;
    }
  }
</script>

<!-- component -->
<div
  class="bg-no-repeat bg-cover bg-center relative
  bg-gradient-to-r from-blue-500 to-purple-500 dark:from-base-100 dark:to-base-100
  "
>
  <div class="absolute sm:inset-0 z-0" />
  <div class="min-h-screen sm:flex sm:flex-row mx-0 justify-center">
    <div class="flex-col flex self-center p-10 sm:max-w-5xl xl:max-w-2xl z-10">
      <div class="self-start hidden lg:flex flex-col text-white">
        <h1 class="mb-3 font-bold text-5xl">
          Hi, Welcome to <span class="font-bold text-5xl"> KubeLab</span>
        </h1>
        <p class="pr-3">Experience Kubernetes Mastery Through Practice.</p>
        <a href="https://natron.io" target="_blank" class="mt-10">
          <span class="text-xs font-semibold leading-6 dark:text-gray-900 text-white"
            >Powered by</span
          >
          {#if $darkTheme === false}
            <img class="h-4 w-auto" src={"/images/natron-dark.png"} alt="Switzerland" />
          {:else}
            <img class="h-4 w-auto" src={"/images/natron.png"} alt="Switzerland" />
          {/if}
        </a>
      </div>
    </div>
    <form
      class="flex justify-center self-center z-10"
      method="POST"
      on:submit|preventDefault={submit}
    >
      <div class="p-12 bg-white dark:bg-neutral mx-auto rounded-2xl w-100">
        <div class="mb-4">
          <ToggleConfetti>
            <div class="btn btn-block btn-ghost normal-case text-xl mb-10" slot="label">
              <img src="/images/kubelab-logo.png" alt="logo" class="w-8 h-8 mr-2" /> KubeLab
            </div>
            <Confetti />
          </ToggleConfetti>
          <h3 class="font-semibold text-2xl">Sign In</h3>
          <p class="text-gray-500">Please sign in to your account.</p>
        </div>
        <div class="space-y-5">
          {#if allowEmailAuth}
            <div class="space-y-2">
              <label class="text-sm font-medium tracking-wide">Email</label>
              <input
                type="text"
                placeholder="your@email.com"
                class="input input-bordered w-full max-w-xs"
                required
                bind:value={user.email}
              />
            </div>
            <div class="space-y-2">
              <label class="mb-5 text-sm font-medium tracking-wide"> Password </label>
              <input
                type="password"
                placeholder="Password"
                class="input input-bordered w-full max-w-xs"
                required
                bind:value={user.password}
              />
            </div>
            <div>
              <button type="submit" class="btn btn-neutral btn-block">
                {#if loading}
                  <span class="loading loading-dots loading-md" /> Loading
                {:else}
                  Sign in
                {/if}
              </button>
            </div>
          {/if}

          {#if allowEmailAuth && allowOAuth2Auth}
            <!-- OAuth Divider -->
            <div class="flex items-center gap-4 py-4">
              <div class="flex-1 border-t border-gray-300 dark:border-gray-600"></div>
              <span class="text-gray-500 text-sm">OR</span>
              <div class="flex-1 border-t border-gray-300 dark:border-gray-600"></div>
            </div>
          {/if}

          {#if allowOAuth2Auth}
            <!-- OAuth Button -->
            <div>
              <button
                type="button"
                class="btn btn-outline btn-block gap-2"
                on:click={handleOAuthLogin}
                disabled={oauthLoading}
              >
                {#if oauthLoading}
                  <span class="loading loading-dots loading-md" /> Loading
                {:else}
                  <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                    <path
                      d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"
                    />
                  </svg>
                  Continue with GitHub
                {/if}
              </button>
            </div>
          {/if}

          {#if !allowEmailAuth && !allowOAuth2Auth}
            <div class="alert alert-error">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="stroke-current shrink-0 h-6 w-6"
                fill="none"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <span>No authentication methods available. Please contact administrator.</span>
            </div>
          {/if}
        </div>
        <div class="pt-5 text-center text-gray-400 text-xs">
          <span>
            Copyright © {new Date().getFullYear()}
            <a
              href="https://natron.io"
              rel=""
              target="_blank"
              title="Natron Tech"
              class="text-blue hover:text-blue-500">Natron Tech</a
            ></span
          >
        </div>
      </div>
    </form>
  </div>
</div>
