<?php include 'header.php'; ?>
    <h2>NIP Checker</h2>
    <form id="nipForm">
        <div class="mb-3">
            <label for="nip" class="form-label">NIP Number</label>
            <input type="text" class="form-control" id="nip" required>
        </div>
        <button type="submit" class="btn btn-primary">Check</button>
        <span class="loading">Loading ...</span>
    </form>
    <div id="nipResult" class="mt-3"></div>

    <?php include 'scripts.php'; ?>
    <script>
        $(document).ready(function() {
            $('#nipForm').submit(function(e) {
                e.preventDefault();
                $('.loading').addClass('show');
                $('#nipResult').empty();
                var nip = $('#nip').val();

                 $.ajax({
                    url: 'http://localhost:20010/nip_checker',
                    type: 'POST',
                    contentType: 'application/json',
                    data: JSON.stringify({
                        nip: nip
                    }),
                     success: function(response) {
                        var resultHtml = '<div class="alert alert-success">';
                        if(response.shortName){
                           resultHtml += '<p><strong>Short Name:</strong> ' + response.shortName + '</p>';
                        }
                        if(response.longName){
                           resultHtml += '<p><strong>Long Name:</strong> ' + response.longName + '</p>';
                        }
                        if(response.taxId){
                           resultHtml += '<p><strong>Tax ID:</strong> ' + response.taxId + '</p>';
                        }

                       if(response.street){
                           resultHtml += '<p><strong>Street:</strong> ' + response.street + '</p>';
                        }
                        if(response.building){
                           resultHtml += '<p><strong>Building:</strong> ' + response.building + '</p>';
                        }
                         if(response.apartment){
                           resultHtml += '<p><strong>Apartment:</strong> ' + response.apartment + '</p>';
                        }
                        if(response.city){
                           resultHtml += '<p><strong>City:</strong> ' + response.city + '</p>';
                        }
                         if(response.province){
                            resultHtml += '<p><strong>Province:</strong> ' + response.province + '</p>';
                        }
                        if(response.zip){
                            resultHtml += '<p><strong>ZIP:</strong> ' + response.zip + '</p>';
                        }

                         resultHtml += '</div>';
                         $('#nipResult').html(resultHtml);
                    },
                    error: function(xhr, status, error) {
                        var errorMessage = 'Error: ' + xhr.status + ' ' + error;
                       if(xhr.responseJSON && xhr.responseJSON.message){
                            errorMessage = 'Error: ' +  xhr.responseJSON.message;
                        }
                         $('#nipResult').html('<div class="alert alert-danger">' + errorMessage  + '</div>');
                    },
                     complete: function(){
                       $('.loading').removeClass('show');
                    }
                });
            });
        });
    </script>

<?php include 'footer.php'; ?>