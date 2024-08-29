// Fetch and display all employees when the page loads
document.addEventListener("DOMContentLoaded", function() {
    fetchEmployees();
    document.getElementById("employeeForm").addEventListener("submit", saveEmployee);
});

// Fetch all employees
function fetchEmployees() {
    fetch("/employees")
        .then(response => response.json())
        .then(data => {
            const tableBody = document.getElementById("employeeTable").querySelector("tbody");
            tableBody.innerHTML = "";  // Clear the table

            data.forEach(emp => {
                const row = document.createElement("tr");

                row.innerHTML = `
                    <td>${emp.firstname}</td>
                    <td>${emp.lastname}</td>
                    <td>${emp.monday}</td>
                    <td>${emp.tuesday}</td>
                    <td>${emp.wednesday}</td>
                    <td>${emp.thursday}</td>
                    <td>${emp.friday}</td>
                    <td>${emp.saturday}</td>
                    <td>${emp.out_of_office ? "Yes" : "No"}</td>
                    <td>${emp.sick ? "Yes" : "No"}</td>
                    <td>
                        <button onclick="editEmployee(${emp.id})">Editer</button>
                        <button class="del" onclick="deleteEmployee(${emp.id})">Supprimer</button>
                    </td>
                `;
                tableBody.appendChild(row);
            });
        })
        .catch(error => console.error("Error fetching employees:", error));
}

// Save or Update an Employee
function saveEmployee(event) {
    event.preventDefault();
    console.log("Form submitted");  // Check if this message appears
    
    const employeeId = document.getElementById("employeeId").value;
    const method = employeeId ? "PUT" : "POST";
    const url = employeeId ? `/employees/${employeeId}` : "/employees";

    const employee = {
        firstname: document.getElementById("firstname").value,
        lastname: document.getElementById("lastname").value,
        monday: document.getElementById("monday").value,
        tuesday: document.getElementById("tuesday").value,
        wednesday: document.getElementById("wednesday").value,
        thursday: document.getElementById("thursday").value,
        friday: document.getElementById("friday").value,
        saturday: document.getElementById("saturday").value,
        out_of_office: document.getElementById("out_of_office").checked,
        sick: document.getElementById("sick").checked
    };

    fetch(url, {
        method: method,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(employee)
    })
    .then(response => {
        if (!response.ok) throw new Error("Failed to save employee.");
        fetchEmployees();  // Refresh the employee list
        document.getElementById("employeeForm").reset();  // Reset the form
        document.getElementById("employeeId").value = '';  
    })
    .catch(error => console.error("Error saving employee:", error));
}

// Edit an Employee
function editEmployee(id) {
    fetch(`/employees/${id}`)
        .then(response => response.json())
        .then(emp => {
            document.getElementById("employeeId").value = emp.id;
            document.getElementById("firstname").value = emp.firstname;
            document.getElementById("lastname").value = emp.lastname;
            document.getElementById("monday").value = emp.monday;
            document.getElementById("tuesday").value = emp.tuesday;
            document.getElementById("wednesday").value = emp.wednesday;
            document.getElementById("thursday").value = emp.thursday;
            document.getElementById("friday").value = emp.friday;
            document.getElementById("saturday").value = emp.saturday;
            document.getElementById("out_of_office").checked = emp.out_of_office;
            document.getElementById("sick").checked = emp.sick;
        })
        .catch(error => console.error("Error fetching employee:", error));
}

// Delete an Employee
function deleteEmployee(id) {
    if (!confirm("Êtes-vous sûr(e) de vouloir supprimer cet employé ?")) return;

    fetch(`/employees/${id}`, { method: "DELETE" })
        .then(response => {
            if (!response.ok) throw new Error("Failed to delete employee.");
            fetchEmployees();  // Refresh the employee list
        })
        .catch(error => console.error("Error deleting employee:", error));
}

document.getElementById("exportBtn").addEventListener("click", function () {
    // Import jsPDF
    const { jsPDF } = window.jspdf;

    // Create a new instance of jsPDF
    const doc = new jsPDF();

    // Define the columns and rows for the table
    const columns = ["Prénom", "Nom", "Lun", "Mar", "Mer", "Jeu", "Ven", "Sam", "Congé", "Arrêt"];
    const rows = [];

    // Get employee data from the table
    const tableRows = document.querySelectorAll("#employeeTable tbody tr");

    tableRows.forEach(row => {
        const rowData = [];
        row.querySelectorAll("td").forEach(cell => {
            rowData.push(cell.innerText);
        });
        rows.push(rowData);
    });

    // Add the table to the PDF
    doc.autoTable({
        head: [columns],
        body: rows,
    });

    // Save the PDF
    doc.save("employee_list.pdf");
});
