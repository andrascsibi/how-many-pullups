angular.module('pullApp')

.controller('ChallengeCtrl', ['$scope', '$rootScope', 'challenge', 'allSets', 'whoami',
  function($scope, $rootScope, challenge, allSets, whoami) {

  $scope.challenge = challenge;
  $rootScope.title = challenge.AccountID + "'s " + challenge.Title + ' challenge';
  $scope.whoami = whoami;

  if (allSets.length === 0) {
    $scope.empty = true;
    return;
  }


  var parseDates = function(sets) {
    return sets
    .filter(function(cur) {
      return cur !== 0;
    })
    .map(function(cur) {
      return {
        ts: Math.round(Date.parse(cur.Date)/1000),
        reps: cur.Reps
      };
    });
  };

  var toCalHeatmap = function(sets) {
    return sets.reduce(function(memo, cur) {
      if (!memo[cur.ts]) {
        memo[cur.ts] = 0;
      }
      memo[cur.ts] += cur.reps;
      return memo;
    }, {});
  };

  var getStats = function(sets) {
    return sets.reduce(function(memo, cur) {
      memo.numSets++;
      memo.totalReps += cur.Reps;
      memo.maxReps = Math.max(memo.maxReps, cur.Reps);
      return memo;
    }, {
      numSets: 0,
      totalReps: 0,
      maxReps: 0,
      minDate: sets[sets.length - 1].Date,
      maxDate: sets[0].Date,
    });
  };

  var groupByDay = function(sets) {
    return sets.reduce(function(memo, cur) {
      var date = new Date(cur.ts*1000).toDateString();
      if (!memo[date]) {
        memo[date] = {sets: [], reps: 0};
      }
      memo[date].sets.push(cur);
      memo[date].reps += cur.reps;
      return memo;
    }, {});
  };

  var sets = parseDates(allSets);
  $scope.stats = getStats(allSets);

  var setsByDay = groupByDay(sets);
  var repsByDay = Object.keys(setsByDay).map(function(k) {
    return setsByDay[k].reps;
  });
  repsByDay.sort(function(l,r) {return l-r;});

  var getPercentile = function(arr, p) {
    arr.sort(function(l,r) {return l-r;});
    return arr[Math.floor(arr.length * p)];
  };

  var getLegend = function(max) {
    var step = max / 4;
    return [1, 2, 3, 4].map(function(e){
      return Math.round(e * step);
    });
  };

  $scope.stats.workDays = repsByDay.length;

  var calSettings = {
    start: new Date($scope.stats.minDate),
    end: new Date(),
    range: 6,
    domain: "week",
    itemName: ['rep', 'reps'],
    tooltip: true,
    data: toCalHeatmap(sets),
    cellSize: 15,
    legendCellSize: 15,
    domainGutter: 10,
  };

  var cal = new CalHeatMap();
  cal.init(angular.extend({
    itemSelector: '#cal-heatmap-punchcard',
    itemNamespace: 'punchcard',
    rowLimit: 24,
    subDomain: "hour",
    legend: [10,20,30,40], // TODO
    legendOrientation: 'vertical',
    legendHorizontalPosition: 'right',
    legendVerticalPosition: 'center',
    label: {position: 'top'},
  }, calSettings));

  var cal2 = new CalHeatMap();
  cal2.init(angular.extend({
    itemSelector: '#cal-heatmap-daily',
    itemNamespace: 'daily',
    rowLimit: 1,
    subDomain: "day",
    subDomainTextFormat: function(date, value) {
      return value;
    },
    legend: getLegend(getPercentile(repsByDay, 0.9)),
    displayLegend: false,
    domainLabelFormat: '',
  }, calSettings));

  $scope.next = function() {
    cal.next();
    cal2.next();
  };

  $scope.previous = function() {
    cal.previous();
    cal2.previous();
  };


}]);
