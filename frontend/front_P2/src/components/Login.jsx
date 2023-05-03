import React, {useState} from 'react'
import { Form, FormGroup, Label, Input, Button} from 'reactstrap'

function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [id, setId] = useState("");

  const handleUsernameChange = (event) => {
    setUsername(event.target.value);
  };
  const handlePasswordChange = (event) => {
    setPassword(event.target.value);
  };
  const handleIdChange = (event) => {
    setId(event.target.value);
  };


  async function loguear(){
    const response = await fetch('http://34.16.139.123:3000/login', {
      method: 'POST',
      body: JSON.stringify({ id, username, password }),
      headers: {
        "Content-Type": "application/json",
      },
    })
    const data = await response.json()
    alert(data.message);
    console.log(data)
    /* if (data.status === "true") {
      const history = useHistory();
      history.push('/');
    } */
  }
  return (
    <>
      <div className='container' style={{margin: "0 auto", maxWidth: '850px', marginTop:"100px"}}>
        <img src='/pngaaa.com-4052093.png' width={200} style={{display: "block", margin: "auto"}}  />
      </div>
      <div className='container' style={{margin: "0 auto", maxWidth: '850px', marginTop:"20px"}}>
        <Form>
        <FormGroup floating>
            <Input
              id="id"
              name="id"
              placeholder="ID Partition"
              type="text"
              style={{ width: "800px" }}
              value={id}
              onChange={handleIdChange}
            />
            <Label for="id">
              ID Partition
            </Label>
          </FormGroup>
          <FormGroup floating>
            <Input
              id="exampleEmail"
              name="user"
              placeholder="Username"
              type="text"
              style={{ width: "800px" }}
              value={username}
              onChange={handleUsernameChange}
            />
            <Label for="exampleEmail">
              Username
            </Label>
          </FormGroup>
          <FormGroup floating>
            <Input
              id="examplePassword"
              name="password"
              placeholder="Password"
              type="password"
              style={{ width: "800px" }}
              value={password}
              onChange={handlePasswordChange}
            />
            <Label for="examplePassword">
              Password
            </Label>
          </FormGroup>
          {' '}
          <Button style={{width: "97%"}} onClick={loguear}>
            Login
          </Button>
        </Form>
      </div>
    </>
    
    
  )
}

export default Login