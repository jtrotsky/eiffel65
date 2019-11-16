import 'core-js-pure/stable';
import 'regenerator-runtime/runtime';
import React, { Component } from 'react';
// import LoginPage from './login/Login';
import Home from './home/Home'
import { BrowserRouter as Router, Redirect, Route, Switch } from 'react-router-dom';
// import { AuthProvider } from './AuthContext'
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles';
import { blue } from './colour/colours';

const theme = createMuiTheme({
  palette: {
    primary: {
      main: blue
    },
  },
  typography: {
    useNextVariants: true,
  }
});

// const getAuthInstance = () => {
//   if (window.gapi === undefined) {
//     console.log("Couldn't get current user, gapi undefined")
//     return;
//   }

//   if (window.gapi.auth2 === undefined) {
//     console.log("gapi loaded, but auth2 not yet loaded")
//     return;
//   }

//   return window.gapi.auth2.getAuthInstance();
// }

// const getCurrentUser = () => {
//   const authInstance = getAuthInstance();

//   if (!authInstance || !authInstance.isSignedIn.get()) {
//     console.log("gapi present, but user not signed in");
//     return;
//   }

//   return authInstance.currentUser.get();
// }

// const isSignedIn = () => {
//   const authInstance = getAuthInstance();

//   return authInstance && authInstance.isSignedIn.get();
// }

class App extends Component {
  state = {
    // isAuthenticated: () => isSignedIn(),
    // getCurrentUser: () => getCurrentUser(),
    refresh: () => this.setState({})
  }

  render() {
    return (
      <Router>
        <MuiThemeProvider theme={theme} >
        <Switch>
            <Route exact path="/">
                <Home />
            </Route>
            <Route path="*">
                <NoMatch />
            </Route>
        </Switch>
        </MuiThemeProvider>
      </Router>
    );
  }
}

function Home() {
    return {Home};
}

function NoMatch() {
    return <h3>404 Not Found</h3>;
}

export default App;