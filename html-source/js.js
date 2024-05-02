document.getElementById('colorButton').addEventListener('click', () => {
    let html = `
      <form id="check-availability-form" action="" method="post" class="needs-validation" novalidate>
        <div class="form-row">
          <div class="col">
            <div class="row" id="reservation-dates-modal">
              <div class="col">
                <input disabled id="start" type="text" name="start" class="form-control" required placeholder="Arival">
              </div>
              <div class="col">
                <input disabled id="end" type="text" name="end" class="form-control" required placeholder="Departure">
              </div>
            </div>
          </div>
        </div>
        <hr>
      </form>
    `;
    attention.custom({msg: html, title: "Check Availability"});
  });