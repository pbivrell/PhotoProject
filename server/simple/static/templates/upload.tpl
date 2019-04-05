<html>
<head>
<script
    src="https://code.jquery.com/jquery-3.3.1.min.js"
    integrity="sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8QtmkMRdAu8="
    crossorigin="anonymous">
</script>
</head>

<h3>Upload Page</h3>

<form>
	<h3>Create or Upload</h3>
    <!-- Switch input type -->
    <input name="projectOrPhotos" id="projectOrPhotos-0" value="project" type="radio">
	<label for="projectOrPhotos-0">Create New Project</label>
	<input name="projectOrPhotos" id="projectOrPhotos-1" value="photos" type="radio" checked="checked">
	<label for="projectOrPhotos-1">Upload Photos</label>
	<input type="hidden" name="token" id="token">
		<h3 >New Project</h3>
			<label for="projectParent">Parent Folder</label>
			<select class="form-control" name="projectParent" id="projectParent">
                {{range $i, $v := .DirectoryStructure }}
                    <option value="{{$v}}">{{$v}}</option> 
			</select>
			<label for="projectName" class="fb-text-label">Project Name</label>
			<input type="text" class="form-control" name="projectName" maxlength="20" id="projectName">
					<input name="checkbox-group-1550036421146[]" id="checkbox-group-1550036421146-0" value="includePhotos" type="checkbox" checked="checked">
					<label for="checkbox-group-1550036421146-0">Include Photos</label>
		<div class="">
			<h3 id="control-9137290">Upload Photos</h3>
		</div>
		<div class="fb-select form-group field-photoParent">
			<label for="photoParent" class="fb-select-label">Parent Folder
				<span class="tooltip-element" tooltip="Select Folder where the images will be uploaded to">?</span>
			</label>
			<select class="form-control" name="photoParent" id="photoParent">
                {{range $i, $v := .DirectoryStructure }}
                    <option value="{{$v}}">{{$v}}</option> 
			</select>
		</div>
		<div class="">
			<h4 id="control-7354177">Source</h4>
		</div>
		<div class="fb-select form-group field-sourceType">
			<label for="sourceType" class="fb-select-label">Source Type
				<span class="tooltip-element" tooltip="Where are the images located?">?</span>
			</label>
			<select class="form-control" name="sourceType" id="sourceType">
				<option value="t1" selected="true" id="sourceType-0">Type1</option>
				<option value="t2" id="sourceType-1">Type2</option>
			</select>
		</div>
		<div class="fb-text form-group field-location">
			<label for="location" class="fb-text-label">Source Location
				<span class="tooltip-element" tooltip="Data required to access the location of the type ie URL">?</span>
			</label>
			<input type="text" class="form-control" name="location" id="location" title="Data required to access the location of the type ie URL">
		</div>
		<div class="fb-button form-group field-remove">
			<button type="button" name="remove" id="remove">Remove Source</button>
		</div>
		<div class="fb-button form-group field-add">
			<button type="button" name="add" id="add">Add Source</button>
		</div>
	</div>
</form>
</html>

<script type="text/javascript">
    function send() {
        if $("#
        var loadData= {
            title1: $("#title1").val(),
            title2: $("#title2").val(),
            link: $("#link").val(),
            description: $("#description").val(),
        }
        
        $.ajax({
            url: '/load',
            type: 'post',
            contentType: 'application/json',
            success: function (data) {
            console.log(data);
            window.location.replace("/display?id="+data);
            
            },
            error: function (msg) {
            document.getElementById("error").innerHTML = msg.responseText;
            document.getElementById("overlay").style.display = "none";
            },
            data: JSON.stringify(loadData)
        });
    }
</script>
