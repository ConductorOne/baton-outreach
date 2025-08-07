## Connector capabilities

1. What resources does the connector sync?
   This connector syncs:
   - Users
   - Profiles
   - Teams

2. Can the connector provision any resources? If so, which ones? 
   The connector can provision:
   - Profile Entitlements
   - Team memberships
   - Accounts

## Connector credentials 

1. What credentials or information are needed to set up the connector? (For example, API key, client ID and secret, domain, etc.)
   - For this connector, an OAuth application on your Outreach account is required.
   
2. For each item in the list above: 

   * How does a user create or look up that credential or info? Please include links to (non-gated) documentation, screenshots (of the UI or of gated docs), or a video of the process. 
      Outreach provides a well-explained documentation on how to set up OAuth: https://developers.outreach.io/api/oauth/
      A simplified version of the steps required to configure it will be:
     1. Log into your Outreach account.
     2. Go to the "Admin Home" panel (hover on the settings icon on the left bar)
     3. Open the "Installed apps" page on the 'Integrations' left section.
     4. Click on the 'Create new app' button.
     5. Click on the 'Create app' button.
     6. Enter the "App name" you desire and click 'Create'.
     7. Once the Details of the new app opens, click on the 'Feature selection' and select 'Outreach API (OAuth)' and click "Add feature".
     8. Copy the OAuth Application ID and Application Secret.
     9. Enter your callback URL.
     10. Select the API Scopes:
         - User: All
         - Teams: All
         - Profiles: All
     11. Save the app and create the release if desired.
      
   * Does the credential need any specific scopes or permissions? If so, list them here. 
     The created App requires the following scopes:
       - User: All
       - Teams: All
       - Profiles: All

   * If applicable: Is the list of scopes or permissions different to sync (read) versus provision (read-write)? If so, list the difference here.
     For read-only:
       - User: Read
       - Teams: Read
       - Profiles: Read

     For read-write:
      - User: All
      - Teams: All
      - Profiles: All

   * What level of access or permissions does the user need in order to create the credentials? (For example, must be a super administrator, must have access to the admin console, etc.)  
      The user should be an admin.