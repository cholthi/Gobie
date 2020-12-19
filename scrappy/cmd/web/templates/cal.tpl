<!DOCTYPE html>
<html lang="en">

<head>

  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
  <meta name="description" content="">
  <meta name="author" content="">

  <title>Agoro Scrapper UI</title>

  <!-- Bootstrap core CSS -->
  <link href="/static/vendor/bootstrap/css/bootstrap.min.css" rel="stylesheet">

  <!-- Custom styles for this template -->
  <link href="/static/css/homepage.css" rel="stylesheet">

</head>

<body>

  <!-- Page Content -->
  <div class="container">
  <div class="row">
    <!-- Sidebar -->
    <div class="col-lg-3">

      <h1 class="my-4">Actions</h1>
      <div class="list-group">
        <a href="/home" class="list-group-item">Scrape Jumia</a>
      </div>
      <div class="list-group">
        <a href="/calculator" class="list-group-item">Price Calculator</a>
      </div>
  </div>
  <!-- End Sidebar -->

  <!-- Main content -->
  <div class="col-lg-9">
    <div class="wrapper">
      <!-- data content -->
      <div class="row">
        <div class="col-2-lg">
          <div class="form-group">
            <label for="formGroupKSh">Jumia Product Price</label>
            <input type="number" class="form-control" id="formGroupKSh" data-bind= "value: JumiaPrice">
        </div>
      </div>
    </div>
    <div class="row">
      <div class="col-lg-4">
        <div>UGX Amount</div>
        <span class="results finalAmount border border-primary" data-bind= "text: UGXProductPriceWithMargin()"></span>
      </div>
      <div class="col-lg-2">
        <div class="form-group">
          <label for="Margin">Margin</label>
        <input type="number" class="form-control" id="Margin" data-bind="value: margin"></input>
      </div>
    </div>
      </div>
    </div>
  </div>
  <!-- End Main Content -->
</div>
  <!-- /.container -->

  <!-- Bootstrap core JavaScript -->
  <script src="/static/vendor/jquery/jquery.min.js"></script>
  <script src="/static/vendor/bootstrap/js/bootstrap.bundle.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/knockout/3.5.0/knockout-min.js"></script>
  <script type="text/javascript">
   $(function(){
     var viewModel = {
          JumiaPrice: ko.observable(0),
          margin: ko.observable(0),

          UGXProductPrice: function(){
            return this.JumiaPrice() * 35
          },

          UGXProductPriceWithMargin: function(){
            return (this.UGXProductPrice() + (this.UGXProductPrice() * this.margin())) + " UGX";
          },

     }
     ko.applyBindings(viewModel);
   })
  </script>

</body>

</html>
