import { client } from "$lib/pocketbase";
import type { PageLoad } from "./$types";

export const load: PageLoad = async () => {
  try {
    const authMethods = await client.collection("users").listAuthMethods();

    return {
      allowEmailAuth: authMethods.usernamePassword ?? true,
      allowOAuth2Auth: (authMethods.authProviders?.length ?? 0) > 0
    };
  } catch (error) {
    console.error("Failed to fetch auth settings:", error);

    // Fallback to defaults if API call fails
    return {
      allowEmailAuth: true,
      allowOAuth2Auth: true
    };
  }
};
