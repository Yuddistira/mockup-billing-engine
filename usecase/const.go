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
						<tbody  id="billing-resp" hx-swap="outerHTML">
						</tbody>
					</table>
				</div>
				<div class="column">
					<h2>Borrower</h2>
					<table style="margin-top: 20px;">
						<tr>
							<td>Status:</th>
							<td id="borrower_status" bgcolor="white" style="color: black; width:max-content;padding: 10px;" >Amount</th>
						</tr>
						<tr>
							<td>Outstanding Amount:</th>
							<td id="borrower_outstanding_amount" bgcolor="white" style="color: black; width:max-content;padding: 10px;" >Amount</th>
						</tr>
					</table>
					<div style="margin-top: 50%;width: 100%;">
						
							
							<h3 style="text-align: center;">Pay amount: <span id="payment_value">50000</span> </h3>
							<div style="width: 100%;text-align: center;">
								<button hx-post="/pay" hx-target="#billing-resp" style="padding: 20px;width: 100px;background-color: aquamarine;">Pay</button>
								<button style="padding: 20px;width: 100px;text-align: center;background-color: firebrick; color: white;">Skip</button>
							</div>							
							
						
						
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
