
export const oktaConfig = {
  issuer: 'https://integrator-2202696.okta.com/oauth2/default',
  clientId: '0oatpv8oadbT62ALK697', 
  redirectUri: window.location.origin + '/login/callback',
  scopes: ['openid', 'profile', 'email'],
  pkce: true,
};
