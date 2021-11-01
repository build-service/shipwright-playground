import React, { useState, useEffect } from 'react';
import { useHistory } from "react-router-dom";
import { PageHeader } from 'antd';
import { NavBar } from '../NavBar/NavBar.js';
import './HowItWorks.css'

const HowItWorks = () => {
  const history = useHistory();
  
  const handleRoute = () =>{ 
    history.push("/");
  }

  return (
    <>
      <NavBar />
      <div className="hiw">
        <PageHeader 
          title="How It Works"
          className="hiw-header"
          onBack={handleRoute}
        />
        <h3 className="hiw-content">Content goes here!</h3>
        <h3 className="hiw-content">Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris nec pretium dolor. Nullam eu massa sed ante facilisis fermentum. Sed sit amet lorem auctor, sollicitudin ipsum at, aliquam enim. Suspendisse sit amet dui urna. Cras at ex neque. Pellentesque mauris enim, lobortis in mi ac, pellentesque vehicula neque. Suspendisse ipsum metus, ultricies vel leo eu, gravida sollicitudin augue. Sed a facilisis purus. Morbi vitae velit eget nulla malesuada volutpat. Quisque ut tempus velit.</h3>
      </div>
    </>
  );
};

export default HowItWorks