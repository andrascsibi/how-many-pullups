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

  var cal = new CalHeatMap();
  cal.init({
    start: new Date($scope.stats.minDate),
    range: 6,
    itemSelector: '#cal-heatmap-punchcard',
    itemNamespace: 'punchcard',
    end: new Date(),
    itemName: ['rep', 'reps'],
    tooltip: true,
    domain: "week",
    rowLimit: 24,
    subDomain: "hour",
    data: toCalHeatmap(sets),
    cellSize: 15,
    legendCellSize: 15,
    legend: [10,20,30,40],
    legendOrientation: 'vertical',
    legendHorizontalPosition: 'right',
    legendVerticalPosition: 'center',
    label: {position: 'top'},
    domainGutter: 10,
  });

  var cal2 = new CalHeatMap();
  cal2.init({
    start: new Date($scope.stats.minDate),
    range: 6,
    itemSelector: '#cal-heatmap-daily',
    itemNamespace: 'daily',
    end: new Date(),
    itemName: ['rep', 'reps'],
    tooltip: true,
    domain: "week",
    rowLimit: 1,
    subDomain: "day",
    data: toCalHeatmap(sets),
    cellSize: 15,
    legendCellSize: 15,
    subDomainTextFormat: function(date, value) {
      return value;
    },
    legend: getLegend(getPercentile(repsByDay, 0.9)),
    displayLegend: false,
    label: {position: 'top'},
    domainGutter: 10,
    domainLabelFormat: '',
  });


}]);
