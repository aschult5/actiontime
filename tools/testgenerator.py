#!/usr/bin/env python3

import argparse
import csv
import logging
import random

from collections import namedtuple
from statistics import mean
from typing import Dict, List


TestCommand = namedtuple('TestCommand', ['cmd', 'action', 'value'])

def gen_action_times(actions: List[str], numadd: int, maxtime: float) -> Dict[str, List[float]]:
    random.seed()

    ret = {a: [random.uniform(0.0001, maxtime) for _ in range(numadd)] for a in actions}
    logging.debug(ret)

    return ret

def gen_test_command(action, times):
    for time in times:
        yield TestCommand('addasync', action, time)

def gen_test_body(actionTimes: Dict[str, List[float]]) -> List[TestCommand]:
    generators = [gen_test_command(action, times) for action,times in actionTimes.items()]

    ret = [val for tup in zip(*generators) for val in tup]
    logging.debug(ret)

    return ret

def gen_test_end(averages: Dict[str, float]) -> List[TestCommand]:
    return []


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Generate test cases.')
    parser.add_argument('action', type=str, nargs='+',
                        help='an action to interleave')
    parser.add_argument('--add', type=int, dest='numadd', action='store',
                        default=1, help='number of calls to add (per action) to interleave')
    parser.add_argument('--maxtime', type=float, dest='maxtime', action='store',
                        default=1E6, help='maximum time an action can take')
    parser.add_argument('--csv', type=str, dest='output', action='store', required=True)
    parser.add_argument('--log', type=str, dest='loglevel', action='store', default='INFO',
                        choices=['DEBUG', 'INFO', 'WARNING', 'ERROR', 'CRITICAL'])

    args = parser.parse_args()

    # Set logging level
    loglevelnum = getattr(logging, args.loglevel.upper(), None)
    if not isinstance(loglevelnum, int):
        raise ValueError('Invalid log level: %s' % args.loglevel)
    logging.basicConfig(level=loglevelnum)

    logging.debug(args)

    # Generate numbers
    actionTimes = gen_action_times(args.action, args.numadd, args.maxtime)
    averages = {a: mean(actionTimes[a]) for a in actionTimes}

    # Interleave numbers
    test = gen_test_body(actionTimes)

    # Check averages
    test += gen_test_end(averages)

    with open(args.output, 'w', newline='') as csvfile:
        testwriter = csv.writer(csvfile, delimiter=',',
                            quotechar='|', quoting=csv.QUOTE_MINIMAL)
        for cmd in test:
            testwriter.writerow(cmd)
