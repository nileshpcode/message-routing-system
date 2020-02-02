from junitparser import TestCase, TestSuite, JUnitXml, Skipped, Error

score_map = {
    'gateway_operations': 5,
    'create_gateway_empty_name': 1,
    'create_gateway_no_ips': 1,
    'create_gateway_same_name': 2,
    'create_gateway_bad_ips': 2,
    'get_gateway_not_exists': 1,
    'route_operations': 5,
    'create_route_empty_prefix': 1,
    'create_route_bad_prefix': 1,
    'create_route_same_prefix': 2,
    'create_route_empty_gateway': 1,
    'create_route_unknown_gateway': 1,
    'get_route_not_exists': 1,
    'route_search': 10,
    'route_search_invalid_number': 1,
}

results = JUnitXml.fromfile('results.xml')
if results._tag == "testsuites":
    cases = [case for cases in results for case in cases]
else:
    cases = list(results)

total = 0
marks = 0
for case in cases:
    score = score_map.get(case.name, 0)
    total += score
    if case.result:
        if case.result._tag != 'failure':
            raise Exception('Invalid junit file')
        continue
    marks += score

score = 100 * marks / total

print('FS_SCORE:%0.2f%%' % score)
