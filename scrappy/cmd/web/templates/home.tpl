<!DOCTYPE html>
<html lang="en">

<head>

  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
  <meta name="description" content="">
  <meta name="author" content="">

  <title>Agoro Scrapper UI</title>

  <!-- Bootstrap core CSS -->
  <link href="/static/vendor/bootstrap/css/bootstrap.min.css" type="text/css" rel="stylesheet">

  <!-- Custom styles for this template -->
  <link href="/static/css/homepage.css" type="text/css" rel="stylesheet">

</head>

<body>

  <!-- Page Content -->
  <div class="container">

    <div class="row">

      <div class="col-lg-3">

        <h1 class="my-4">Actions</h1>
        <div class="list-group">
          <a href="/home" class="list-group-item">Scrape Jumia</a>
        </div>
        <div class="list-group">
          <a href="/calculator" class="list-group-item">Price Calculator</a>
        </div>
      </div>
      <!-- /.col-lg-3 -->

      <div class="col-lg-9">
      <div class="wrapper">
        <div class="row">
        <div class="col-lg-8 scrape-form">
        <div data-bind="hidden: scrappedFinished">
          <h3>Start Scrapping</h3>
          <form data-bind="submit: onScrapeSubmit">
  <div class="form-group">
    <label for="formGroupExampleInput">Jumia Country</label>
    <input type="text" class="form-control" id="formGroupExampleInput" data-bind= "value: host" placeholder="">
  </div>
  <div class="form-group">
    <label for="formGroupExampleInput2">Category To Scrape</label>
    <input type="text" class="form-control" id="formGroupExampleInput2" data-bind= "value: categorytoScrape" placeholder="">
  </div>
  <div class="form-group">
    <label for="formGroupExampleInput2">Replace Text</label>
    <input type="text" class="form-control" id="formGroupExampleInput2" data-bind= "value: replace" placeholder="">
  </div>
  <div class="form-group">
    <label for="formGroupExampleInputfile">Scrape File</label>
    <input type="text" class="form-control" id="formGroupExampleInputfile" data-bind= "value: file" placeholder="./scrapped.json">
  </div>
  <div class="form-group">
    <label for="formGroupExampleInput2">Number of Products to Scrape</label>
    <input type="number" class="form-control" id="formGroupExampleInput2" data-bind= "value: noRequests">
  </div>
  <button type="submit" class="btn btn-primary mb-2" data-bind="disable: disableButton">Scrape Products</button>
</form>
</div>

        <!-- make it invisible only javascript -->
        <div class="upload-form" data-bind="visible: scrappedFinished">
         <h3> Upload products Cscart Shop</h3>
          <form data-bind="submit: onUploadSubmit">
          <div class="form-group">
            <label for="vendor">Vendor ID</label>
            <input type="number" class="form-control" id="vendor" data-bind="value: vendorID">
          </div>
          <div class="form-group">
            <label for="Category">Category ID to place products under</label>
            <input type="number" class="form-control" id="Category" data-bind="value: productCategory">
          </div>
          <div class="form-group">
            <label for="ForexRate">Exchange rate to UGX</label>
            <input type="number" class="form-control" id="ForexRate" data-bind="value: rate">
          </div>
          <div class="form-group">
            <label for="margin">Margin as a percentage of product Price</label>
            <input type="input" class="form-control" id="margin" data-bind="value: margin">
          </div>
          <div class="form-group">
            <input type="hidden" class="form-control" data-bind= "value: file" >
          </div>
          <button type="submit" class="btn btn-primary mb-2">Upload to Shop</button>
        </form>
        </div>
</div>
        <!-- make it invisible only show javascript -->

        <div class="col-lg-4 border" id="results">
         <h4>Results</h4>
          <div id="scrapeResults" data-bind="visible: scrappedFinished">
          <div data-bind="visible: scrappedSuccess">
            <div class="p-3 mb-2 bg-success text-white">Scrapped <span id="number" data-bind="text: noScrapped"></span></div>
          </div>
          <div data-bind="hidden: scrappedSuccess">
            <div class="p-3 mb-2 bg-danger text-white">Scrapped Failed, Please refresh and try again</div>
          </div>
        </div>
         <div data-bind="visible: showScrapeLoader"><img src="/static/assets/img/load.gif" width="80" height="60"/></div>

        <div id="uploadResults" data-bind="visible: uploadFinished">
          <div data-bind="visible: uploadSuccess">
            <div class="p-3 mb-2 bg-primary text-white">Uploaded Sucessefully to Shop</div>
          </div>
          <div data-bind="hidden: uploadSuccess">
            <div class="p-3 mb-2 bg-danger text-white">Upload Failed, Please refresh and try again</div>
        </div>

        </div>
        <div data-bind="visible: showUploadLoader"><img src="/static/assets/img/load2.gif" width="80" height="60"/></div>
  </div>

        <!-- /.row -->

      </div>
      <!-- /.col-lg-9 -->

    </div>
    <!-- /.row -->
    </div>
  </div>
</div>
  <!-- /.container -->

  <!-- Bootstrap core JavaScript -->
  <script src="/static/vendor/jquery/jquery.min.js" type="text/javascript"></script>
  <script src="/static/vendor/bootstrap/js/bootstrap.bundle.min.js" type="text/javascript"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/knockout/3.5.0/knockout-min.js" type="text/javascript"></script>
  <script type="text/javascript">
   $(function(){
     var viewModel = {
       host : ko.observable(""),
       categorytoScrape : ko.observable(""),
       file : ko.observable(""),
       replace : ko.observable(""),
       noRequests : ko.observable(""),
       productCategory : ko.observable(""),
       vendorID : ko.observable(""),
       margin : ko.observable(1),
       rate:    ko.observable(1.1),
       scrappedFinished : ko.observable(false),
       uploadFinished : ko.observable(false),
       noScrapped : ko.observable(0),
       noUploaded : ko.observable(0),
       showScrapeLoader : ko.observable(false),
       showUploadLoader : ko.observable(false),
       uploadSuccess : ko.observable(false),
       scrappedSuccess : ko.observable(false),
       disableButton : ko.observable(false),

       scrape : function() {
         var data = {
                      host: this.host(),
                      category: this.categorytoScrape(),
                      request_no: this.noRequests(),
                      replace: this.replace(),
                      file: this.file(),
         }
          console.log(data)
         //var jsond = ko.toJSON(data)
         var that = this

         $.post("/ajax/scrape", data, function(response){
            console.log(response)
           that.showScrapeLoader(false)
           if (response.success == true) {
              that.scrappedSuccess(true)
             that.scrappedFinished(true)
             that.noScrapped(response.number)
             that.file(response.file)
              return
           }
            that.scrappedFinished(true)
            that.scrappedSuccess(false)
         })
       },
       upload : function() {
         var data = {
                      category: this.productCategory(),
                      vendor: this.vendorID(),
                      rate:  this.rate(),
                      margin: this.margin(),
                      file: this.file(),
         }
         //var jsond = ko.toJSON(data)
         var that = this

         $.post("/ajax/upload", data, function(response){
            that.showUploadLoader(false)
            console.log(response)
           if (response.success == true) {
             that.uploadSuccess(true)
             that.uploadFinished(true)
             that.noUploaded(response.number)
              return
           }
            that.uploadFinished(true)
            that.uploadSuccess(false)
         })
       },
       onScrapeSubmit : function(formElement) {
         this.showScrapeLoader(true)
         this.scrape()
         this.disableButton(true)
       },
       onUploadSubmit : function(formElement) {
         this.showUploadLoader(true)
         this.upload()
       },
     }
     ko.applyBindings(viewModel);
   })
  </script>

</body>

</html>
