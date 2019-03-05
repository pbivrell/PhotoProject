# Photo Project
 
  The photo project is intended to wrap any image storage space in an elegant, easy to navigated, programatically generated web page. This web app is designed to be highly configureable and easy to deploy.

### Endpoints
Creat

### Configuration 
### JSON data:
=====================
"configData": {
    "ip": "";
    "port": "";
}

config.json | configData
    ip, port: Self explanitory

"setupData": {
    "root-link": "";
    "storage-type": "";
}

storage-link.json | setupData
    root-link: Link to storage page
    storage-type: Enumeration of storage types


"pageData": {
    "main-links": [];
    "sub-links": [];
    "link-titles": [];
    "page-title": "";
    "page-subtitle":"";
    "page-description":"";
}

photo-page.json | pageData
    main-links: Large image
    sub-links: Tiny image
    link-titles: N/A
    page-title, page-subtitle, page-description: self explainitory

driver-page.json | pageData
    main-links: First images from each of the sub-directories
    sub-links: Link to the sub-directory page
    link-titles: Text describing whats in the sub-directory
    page-title: N/A
    page-subtitle: N/A
    page-description: N/A



=====================
Storage Structure:
=====================
- Root
    - driver-page.json
    - trip1/
        - driver-page.json
        - Adventure1/
            - photo-page.json
            - jpges+
        - Adventure2/
            - photo-page.json
            - jpges+
    - trip2/
        - photo-page.json
        - jpges+

=====================
Web Server Activity:
=====================
1. Startup
    A. Configure:
       The webapp first attempts to configure the webserver by loading the configuration file specified by the command line.
       *The default configuration file is configuration.json*
    B. Fetch Root of Project:
       The webapp then attempts to load the main page by reading the storage configuration file specified by the configuration file. This will give you the path to the read storage for driver-page.json which contains the json data to generate the home page. If the storage configuration file is not specified the person starting the server will be notified via the command line and all access to the webpage will be redirected to a configuration page. Here the server runner will be able to specify the read and write storage types.
    C. Create:
