/** biome-ignore-all lint/suspicious/noDocumentCookie: CookieStore API is not widely available */
import { redirect } from '@remix-run/react';

const LOGGED_IN_NAME = '_me_logged_in';

// After successful authentication, create an additional non-httpOnly cookie
// to indicate the user is logged in. This cookie allows us to check the user's
// logged-in status without making an API request. Actual authentication for API
// requests will still rely on the secure httpOnly cookie.
const LOGGED_IN_COOKIE = `${LOGGED_IN_NAME}=true; Path=/; SameSite=Lax; Secure`;

const EXPIRE_LOGGED_IN = `${LOGGED_IN_NAME}=; Max-Age=0; Path=/; SameSite=Lax; Secure`;

const hasSession = () =>
	document.cookie
		.split(';')
		.some((c) => c.trim().startsWith(`${LOGGED_IN_NAME}=`));

const expireSession = (shouldRedirect = true) => {
	// Expire the logged-in cookie
	document.cookie = EXPIRE_LOGGED_IN;

	// Check if we should redirect the user to the login page
	// and whether to include the current path as a redirect parameter
	if (shouldRedirect) {
		const currentPath = window.location.pathname;
		if (currentPath !== '/' && currentPath !== '/login') {
			return redirect(`/login?redirect=${currentPath}`);
		}

		return redirect('/login');
	}
};

export { EXPIRE_LOGGED_IN, expireSession, hasSession, LOGGED_IN_COOKIE };
