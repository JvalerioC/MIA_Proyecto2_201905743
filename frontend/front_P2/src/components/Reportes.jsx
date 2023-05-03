import React, { useState, useEffect } from 'react'
import { Form, FormGroup, Label, Input, Button} from 'reactstrap'

function Reportes() {
  const [selectedRep, setSelectecRep] = useState("");
  const [repContent, setRepContent  ] = useState([]);
  const [fileType, setFileType] = useState("");

  const handleRepSelect = async (rep, report) => {
    const response = await fetch(`http://34.16.139.123:3000/reports/${rep}`);
    const data = await response.json();
    console.log(data.message)
    const decodedContent = atob(data.message)
    setSelectecRep(report);

    const nombreArchivo = report.split('/').pop();
    const extension = nombreArchivo.split('.').pop();
    var fileType = ""
    switch (extension){
      case "pdf":
        setFileType("application/pdf")
      case "png":
        setFileType("image/png")
      case "jpg":
      case "jpeg":
        setFileType("image/jpeg")
      case "txt":
        setFileType("text/plain")
    }
    downloadFile(decodedContent, nombreArchivo)

  };

  function downloadFile(content, fileName) {
    const byteCharacters = content;
    const byteNumbers = new Array(byteCharacters.length);
    for (let i = 0; i < byteCharacters.length; i++) {
      byteNumbers[i] = byteCharacters.charCodeAt(i);
    }
    const byteArray = new Uint8Array(byteNumbers);
    const blob = new Blob([byteArray], { type: fileType });
    //const blob = new Blob([content], { type: fileType });
    const link = document.createElement("a");
    link.href = window.URL.createObjectURL(blob);
    link.download = fileName;
    link.click();
  }

  const contenidoReporte = async () => {
    const response = await fetch(`http://34.16.139.123:3000/reports`);
    const data = await response.json();
    //const data = await response.text()
    console.log(data)
    setRepContent(data.message);
  };
  useEffect(() => {
    contenidoReporte();
  },[])

  return (
    <>
      <div className='container' style={{margin: "0 auto", maxWidth: '700px', marginTop:"10px"}}>
        <h1 style={{textAlign: "center", marginTop: "10px", color: "grey", marginBottom:"50px"}}>Reportes</h1>
      <ul >
        {repContent.map((report, index) => (
          <li key={index} onClick={() => handleRepSelect(index, report)}>
            {report}
          </li>
        ))}
      </ul>
      </div>
    </>
  )
}

export default Reportes