import React, { useState, useRef } from 'react'
import { Form, FormGroup, Label, Input, Button} from 'reactstrap'

function Inicio() {
  const [selectedFile, setSelectecFile] = useState(null);
  const [fileContent, setFileContent  ] = useState('');
  const [contentConsole, setContentConsole ] = useState("")
  const inputFile = useRef(null);

  const handleFileSelect = (event) => {
    const file = event.target.files[0];
    setSelectecFile(file);
    const reader = new FileReader();
    reader.onload = (e) => {
      setFileContent(e.target.result);
    };
    reader.readAsText(file);
  };

  const handleChange = (event) => {
    setFileContent(event.target.value);
  };
  async function ejecutar(){
    const response = await fetch('http://34.16.139.123:3000/execute', {
      method: 'POST',
      body: JSON.stringify({ fileContent }),
      headers: {
        "Content-Type": "application/json",
      },
    })
    const data = await response.json()
    console.log(data.message)
    setContentConsole(data.message)
  }
  const limpiar = () => {
    //setResult(result.slice(0, result.length -1))  
    setFileContent("")
  }

  return (
    <>
      <div className='container' style={{margin: "0 auto", maxWidth: '850px', marginTop:"10px"}}>
        <FormGroup>
          <Label for="exampleFile">
            Abrir Archivo 
          </Label>
          <Input
            id="exampleFile"
            name="file"
            type="file"
            onChange={handleFileSelect}
            ref={inputFile}
          />
        </FormGroup>
        <Form>
          <FormGroup>
            <Label for="exampleTexti">
              Entrada
            </Label>
            <Input
              id="exampleTexti"
              name="texti"
              type="textarea"
              rows={12}
              onChange={handleChange}
              value={fileContent}
            />
          </FormGroup>
          <Button style={{width: "49%", marginRight: "10px"}} onClick={ejecutar}>
            Ejecutar
          </Button>
          <Button style={{width: "48%", marginLeft: "10px"}} onClick={limpiar}>
            Limpiar
          </Button>
          <br/><br/>
          <FormGroup>
            <Label for="exampleTexto">
              Salida
            </Label>
            <Input
              id="exampleTexto"
              name="texto"
              type="textarea"
              rows={12}
              value={contentConsole}
              readOnly
            />
          </FormGroup>
          
        </Form>
      </div>
    </>
  )
}

export default Inicio;
