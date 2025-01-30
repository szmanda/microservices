<?php include 'header.php'; ?>
    <h2>Submit Print Job</h2>
    <form id="printForm">
        <div class="mb-3">
            <label for="documentName" class="form-label">Document Name</label>
            <input type="text" class="form-control" id="documentName" required>
        </div>
        <div class="mb-3">
            <label for="paperSize" class="form-label">Paper Size</label>
            <select class="form-select" id="paperSize" required>
                <option value="A4">A4</option>
                <option value="A3">A3</option>
                <option value="Letter">Letter</option>
            </select>
        </div>
        <div class="mb-3">
            <label for="orientation" class="form-label">Orientation</label>
            <select class="form-select" id="orientation" required>
                <option value="portrait">Portrait</option>
                <option value="landscape">Landscape</option>
            </select>
        </div>
        <div class="mb-3">
            <label for="copies" class="form-label">Copies</label>
            <input type="number" class="form-control" id="copies" value="1" min="1" required>
        </div>
        <div class="mb-3">
            <label for="printers" class="form-label">Printers</label>
            <select class="form-select" id="printers" multiple required>
               <option value="printer1">Printer 1</option>
               <option value="printer2">Printer 2</option>
               <option value="printer3">Printer 3</option>
             </select>
             <div id="printersHelp" class="form-text">Hold down the Ctrl (windows) or Command (Mac) key to select multiple printers.</div>
        </div>
        <button type="submit" class="btn btn-primary">Submit</button>
        <span class="loading">Loading ...</span>

    </form>
    <div id="printResult" class="mt-3"></div>
    <?php include 'scripts.php'; ?>
    <script>
        $(document).ready(function() {
            $('#printForm').submit(function(e) {
                e.preventDefault();
                $('.loading').addClass('show');
                $('#printResult').empty(); // Clear previous messages
                var document_name = $('#documentName').val();
                var paper_size = $('#paperSize').val();
                var orientation = $('#orientation').val();
                var copies = parseInt($('#copies').val());
                var printers = $('#printers').val();

                $.ajax({
                    url: 'http://localhost:20010/print',
                    type: 'POST',
                    contentType: 'application/json',
                    data: JSON.stringify({
                        document_name: document_name,
                        paper_size: paper_size,
                        orientation: orientation,
                        copies: copies,
                        printers: printers
                    }),
                    success: function(response) {
                        $('#printResult').html('<div class="alert alert-success">' + response.message + '</div>');
                    },
                    error: function(xhr, status, error) {
                        var errorMessage = 'Error: ' + xhr.status + ' ' + error;
                        if(xhr.responseJSON && xhr.responseJSON.message){
                            errorMessage = 'Error: ' +  xhr.responseJSON.message;
                        }

                       $('#printResult').html('<div class="alert alert-danger">' + errorMessage  + '</div>');
                    },
                    complete: function(){
                       $('.loading').removeClass('show');
                    }
                });
            });
        });
    </script>
<?php include 'footer.php'; ?>