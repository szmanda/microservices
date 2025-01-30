<?php include 'header.php'; ?>
    <h2>Check Print Job Status</h2>
    <form id="statusForm">
         <button type="submit" class="btn btn-primary">Get Status</button>
        <span class="loading">Loading ...</span>
    </form>
    <div id="statusResult" class="mt-3"></div>
    <?php include 'scripts.php'; ?>
    <script>
        $(document).ready(function() {
            $('#statusForm').submit(function(e) {
                e.preventDefault();
                $('.loading').addClass('show');
                $('#statusResult').empty(); // Clear previous messages

                $.ajax({
                    url: 'http://localhost:20010/print/status',
                    type: 'GET',
                    dataType: 'json', // Expect JSON response
                    success: function(response) {
                       console.log(response);
                        if(response && response.message) {
                           $('#statusResult').html('<div class="alert alert-success">Status: ' + response.message + '</div>');
                        } else {
                           $('#statusResult').html('<div class="alert alert-warning">Invalid response format received.</div>');
                        }
                    },
                    error: function(xhr, status, error) {
                        var errorMessage = 'Error: ' + xhr.status + ' ' + error;
                         if(xhr.responseJSON && xhr.responseJSON.message){
                            errorMessage = 'Error: ' +  xhr.responseJSON.message;
                         }
                        $('#statusResult').html('<div class="alert alert-danger">' + errorMessage + '</div>');
                    },
                    complete: function(){
                       $('.loading').removeClass('show');
                    }
                });
            });
        });
    </script>
<?php include 'footer.php'; ?>