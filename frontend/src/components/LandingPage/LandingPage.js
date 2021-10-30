import React from 'react';
import background from '../../assets/images/shipit.jpg';
import { NavBar } from '../NavBar/NavBar.js';
import { Button } from 'antd';
import './LandingPage.css'

export const LandingPage = () => (
  <>
    <NavBar />
    <img src={background} className="background" alt="a container ship at sea" />
    <div className="float">
      <div className="blurb">A super catchy sentence that will capture your Attention. <br/>The Future of Building. </div>
    </div>
    <div className="float text-center">
      <Button shape="round" size={"large"} className="start-button" href="#form">Get Started</Button>
    </div>
    <div className="image-credit">Image Courtesy of Martin Dörsch, stocksnap.io</div>
  </>
)