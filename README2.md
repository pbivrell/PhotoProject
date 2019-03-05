# Photo Project
 
  The photo project is intended to wrap any image storage space in an elegant, easy to navigate, programatically generated web page. This web app is designed to be highly configurable and easy to deploy.

### Endpoints
Below Describes the Endpoints served by the photo project server.
___
#### Upload [Page]:
* __URL:__ `/upload`
* __HTTP Method__: `GET`
* __Description:__
This endpoint servers the HTML page that provides inputs for the JSON data need for the [`/uploadProject`]() and [`/uploadImages`]() endpoints. 

* __Url params__: 
   Required: None
   _Optional_: 
   `token=`one-time-token --- see [`/admit`]()
* __Data Params:__
None
* __Success:__
Status Code: 200 OK
Content: [upload.html]()
* __Error:__
Status Code: 404 Not Found
___
#### Upload Project :
* __URL:__ `/uploadProject`
* __HTTP Method__: `POST`
* __Description:__
This endpoint creates a new project. A project is a directory for image files to be stored under.
* __Url params__: 
   Required: None
   _Optional_: None
   
* __Data Params:__
```json
{
	"parent": "the parent directory of the new project. Either 'root' or an existing directory",
	"name":"name of new project. Name must be unique within the parent directory."
	"token": "one-time-token this will authorize the creation of this project. This is only required if authorization is required."
}
```
* __Success:__
Status Code: 303 See Other
Headers: 
Location: [`/load?jid=XXX`]()

* __Error:__
Status Code: 500 Internal Server Error
___
#### Upload Images :
* __URL:__ `/uploadImages`
* __HTTP Method__: `POST`
* __Description:__
This endpoint fetches images from the provided storage volume and uploads photos to the desired project. 
* __Url params__: 
   Required: None
   _Optional_: None
   
* __Data Params:__
```json
{
	"parent": "the parent project of the images. Must be an existing project. A project can either contain images or sub-projects images can not be uploaded to projects which contain sub-projects.",
	"source": [
		"type": "one of the valid read storage types. No type implies URL",
		"url": "path to the image or images"
		"permissions": true or false // Are permissions required to access read storage
	]
	"token": "one-time-token this will authorize the creation of this project. This is only required if authorization is required."
}
```
* __Success:__
Status Code: 303 See Other
Headers: 
Location: [`/load?jid=XXX`]()
_OR_
Status Code: 303 See Other
Headers:
Location: Read storage permission page

* __Error:__
Status Code: 500 Internal Server Error
___
#### Admit [Page]:
* __URL:__ `/admit`
* __HTTP Method__: `GET`
* __Description:__
This endpoint generates a one-time-token that can be used to authorize access to the [`/upload`]() endpoint. Access to this page requires a password set at configuration time.
* __Url params__: 
   Required: 
   `code=`Password for admit page 
   _Optional_: 
   `exp=`Expiration time in mins for token. Default is 1440 mins.
   `uses=`Number of times this key can be used. Default is 1.
 
   
* __Data Params:__
None
* __Success:__
Status Code: 200 OK
Content:
URL to page needing authorization including URL encoded token

* __Error:__
Status Code: 500 Internal Server Error
___

#### Loading [Page]:
* __URL:__ `/loading`
* __HTTP Method__: `GET`
* __Description:__
This endpoint serves the loading webpage which displays the progress of the uploading images. The progress is retrieved with a GET request to [`/progress`]()
* __Url params__: 
   Required: 
   `jid=` Job id for the 
   _Optional_: None
   
* __Data Params:__
None
* __Success:__
Status Code: 200 OK
Content: [loading.html]()

* __Error:__
None
___

#### Progress:
* __URL:__ `/progress`
* __HTTP Method__: `GET`
* __Description:__
This endpoint returns text representing the progress of images being uploaded.
 
* __Url params__: 
   Required: 
   `jid=` Job id for the 
   _Optional_: None
   
* __Data Params:__
None
* __Success:__
Status Code: 200 OK
Content: 
Text representing the progress of images being uploaded

* __Error:__
None
___

#### Display [Page]:
* __URL:__ `/{path:*}`
* __HTTP Method__: `GET`
* __Description:__
This endpoint generates and returns the webpage for the project or photo page found at the path specified. This endpoint matches all other endpoints however the functional endpoint will take precedence. 
 
* __Url params__: 
   Required: None
   _Optional_: None
   
* __Data Params:__
None

* __Success:__
Status Code: 200 OK
Content: [photoPage.tpl]() or [projectPage.tpl]()

* __Error:__
Status Code: 404 Not Found
___

#### Configure [Page]:
* __URL:__ `/configure`
* __HTTP Method__: `GET`
* __Description:__
This endpoint generates the webpage that lists the not configured storage volumes. This page will include instructions on how to configure these storage volumes. All other endpoints will redirect to here when a volume has not be configured.
 
* __Url params__: 
   Required: None
   _Optional_: None
   
* __Data Params:__
None

* __Success:__
Status Code: 200 OK
Content: [configured.html]()

* __Error:__
None
___

