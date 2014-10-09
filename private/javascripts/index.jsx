/** @jsx React.DOM */

function roundToDecimal(num, decimals) {
  var multiplier = Math.pow(10, decimals);
  return Math.round(num * multiplier) / multiplier;
};

function formatPercent(num) {
  return roundToDecimal(num * 100, 2);
};

function fetchCohorts(numCohorts, cohortDuration) {
  return $.get("/data", {
    numCohorts: numCohorts,
    cohortDuration: cohortDuration
  });
};

function formatDate(time) {
  var date = new Date(time);
  return (date.getUTCMonth()+1) + "/" + date.getUTCDate();
}

var DetailsComponent = React.createClass({
  render : function() {
    var orderers = this.props.orderers;
    var firsts = this.props.firsts;
    var totalUsers = this.props.users;
    return (
      <div className="cohorts-table-column">
        <div>{formatPercent(orderers / totalUsers)}% orderers ({orderers})</div>
        <div>{formatPercent(firsts / totalUsers)}% 1st time ({firsts})</div>
      </div>
    )
  }
});

var CohortComponent = React.createClass({
  render : function() {
    var self = this;

    var details = _.map(this.props.details, function(detail) {
      return (
        <DetailsComponent key={detail.start} orderers={detail.orderers} firsts={detail.firsts} users={self.props.users} />
      )
    });

    return (
      <div className="cohorts-table-row">
        <div className="cohorts-table-column">
          {formatDate(this.props.start)}-{formatDate(this.props.end)}
        </div>
        <div className="cohorts-table-column">
          {this.props.users} users
        </div>
        {details}
      </div>
    )
  }
});

var CohortsTableComponent = React.createClass({
  render : function() {
    var dayHeaders = [];
    for (var i = 0; i < this.props.numCohorts; i++) {
      dayHeaders.push(
        <div className="cohorts-table-header" key={i}>{i * this.props.cohortDuration}-{(i+1) * this.props.cohortDuration} days</div>
      )
    }

    var cohorts = _.map(this.props.cohorts, function(cohort) {
      return (
        <CohortComponent key={cohort.start} start={cohort.start} end={cohort.end} users={cohort.users} details={cohort.details} />
      )
    });

    return (
      <div className="cohorts-table">
        <div className="cohorts-table-row">
          <div className="cohorts-table-header">Cohort</div>
          <div className="cohorts-table-header">Users</div>
          {dayHeaders}
        </div>
        {cohorts.reverse()}
      </div>
    );
  }
});

var InterfaceComponent = React.createClass({
  getInitialState: function() {
    return {
      numCohorts: 8,
      cohortDuration: 7,
      cohorts: null,
      hasError: false,
      fetching: false,
      displayedNumCohorts: null,
      displayedCohortDuration: null
    }
  },
  render : function() {
    var result;
    if (this.state.hasError) {
      result = <div className="error">An error occured, please try again</div>
    } else if (this.state.cohorts) {
      result = <CohortsTableComponent
            numCohorts={this.state.displayedNumCohorts}
            cohortDuration={this.state.displayedCohortDuration}
            cohorts={this.state.cohorts} />
    } else if (this.state.fetching) {
      result = <div className="alert">Loading!</div>
    }

    var fetchButton = !this.state.fetching ?
      <button className="btn btn-default btn-sm" onClick={this.fetch}>Fetch</button> :
      <button className="btn btn-default btn-sm" disabled="disabled">Fetching</button>;

    return (
      <div>
        <div className="controls-container">
          <div className="control">
            <label>Cohorts</label>
            <input type="number" min="1" step="1" defaultValue="8" placeholder="8" onChange={this.changeWeeks} />
          </div>
          <div className="control">
            <label>Cohort duration</label>
            <input type="number" min="1" step="1" defaultValue="7" placeholder="7" onChange={this.changeCohortDuration} />
            days
          </div>
          <div className="control">
            {fetchButton}
          </div>
        </div>
        {result}
      </div>
    );
  },
  changeWeeks: function(e) {
    this.setState({numCohorts: e.target.value || 8});
  },
  changeCohortDuration: function(e) {
    this.setState({cohortDuration: e.target.value || 7});
  },
  fetch: function() {
  var self = this;

    this.setState({cohorts: null, hasError: false, fetching: true});
  fetchCohorts(this.state.numCohorts, this.state.cohortDuration).done(function(data) {
    self.setState({
      cohorts: data,
      hasError: false,
      fetching: false,
      displayedNumCohorts: self.state.numCohorts,
      displayedCohortDuration: self.state.cohortDuration});
  }).fail(function() {
    self.setState({cohorts: null, hasError: true, fetching: false});
  });
  }
});

$(function() {
  React.renderComponent(
    <InterfaceComponent />,
    document.getElementById('main-container')
  );
});