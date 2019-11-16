import React, { Component } from 'react';
import { Redirect, Route, Switch } from 'react-router-dom';
import MarketClient from './market/MarketClient';
// import Snackbar from '@material-ui/core/Snackbar';
import IconButton from '@material-ui/core/IconButton';
import CloseIcon from '@material-ui/icons/Close';
import { withStyles } from '@material-ui/core/styles';

// function Home(props) {
//     return (
//       <AuthConsumer>
//         {authContext => <HomePage {...props} authContext={authContext} /> }
//       </AuthConsumer>
//     );
//   }

const styles = theme => ({
    close: {
        padding: 4,
    },
});

class HomePage extends Component {
    state = {}

    // displaySnackbar(message) {
    //     this.setState({
    //         snackbar: {
    //             open: true,
    //             message: message
    //         }
    //     });
    // }

    // handleSnackbarClose() {
    //     this.setState({
    //         snackbar: {
    //             open: false,
    //             message: null
    //         }
    //     });
    // }

    constructor(props) {
        super(props);

        this.marketClient = new MarketClient();

        this.getListings();
    }

    getListings() {
        this.marketClient.getListings().then(listings => {
                console.log(listings);
                let listin = users.map(user => ({
                    name: user.FIRST_NAME,
                    value: user.USER_ID
                }));
                let selectedUsers = new Set();
                reportUsers.forEach(user => selectedUsers.add(user.value));

                this.setReportContext({reportUsers, selectedUsers});

                this.generateReport(createDefaultReportData(reportUsers.map(user => user.value)));
        }).catch((err) => {
            console.log(err);
        });
    }


    render() {
        const { match, classes } = this.props;

        return (
            <div style={{height: "64px"}} />
            <Route path={`/home`} render={('test')} />
        );
    }
}

export default withStyles(styles)(Home);