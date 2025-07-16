package usecase

const SimulationResultPage string = `
<hr>
			<h1>Simulation</h1>
			<div class="row">
				<div class="column">
					<table border="1" style="margin-top: 20px;">
						<thead>
							<tr>
								<th>Billing Schedule</th>
								<th>Amount</th>
								<th>Payment Status</th>
							</tr>
						</thead>
						<tbody id="billing-resp">
							{{range .Billings}}
								<tr>
									<td>{{.PaymentIdx}}</td>
									<td>{{.Amount}}</td>
									<td>{{.PayStatus}}</td>
								</tr>
							{{end}}
						</tbody>
					</table>
				</div>
				<div class="column">
					<h2>Borrower</h2><h2> {{.Finish}}</h2>


					<div style="width: 100%;text-align: left;">
						<td>Status:</td>
					</th>
					<td id="borrower_status" hx-swap-oob="true">{{.Status}}</td>
					</div>

					<td>Outstanding Amount:</td>
					<td id="borrower_outstanding_amount" hx-swap-oob="true">Rp {{.OutstandingAmount}}</td>


					<h3 style="text-align: center;">{{.BillSchedule}} Pay amount: {{.Interest}}</h3>
					<div id="interest-container" data-interest="{{.Interest}}" data-billingid="{{.BillingID}}"></div>
					<div style="width: 100%;text-align: center;">
						<button 
							hx-post="/pay" 
							hx-target="#response" 
							hx-vals='js:{
								"interest": document.getElementById("interest-container").dataset.interest,
								"billing_id": document.getElementById("interest-container").dataset.billingid
							}' 
							style="padding: 20px;width: 100px;background-color: aquamarine;">
							Pay
						</button>
						<button
							hx-post="/skip" 
							hx-target="#response" 
							hx-vals='js:{
								"billing_id": document.getElementById("interest-container").dataset.billingid
							}'
							style="padding: 20px;width: 100px;text-align: center;background-color: firebrick; color: white;">
							Skip
						</button>
					</div>


				</div>
			</div>`

const NewRowSimulationTable string = `
<tr>
    <td>test1</td>
    <td>test2</td>
    <td>test3</td>
</tr>
`
