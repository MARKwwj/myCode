import pytest


def inc(x):
    return x + 1


def test_answer():
    assert inc(3) == 5


@pytest.fixture(scope='function')
def test_print():
    assert 1 == 1


if __name__ == '__main__':
    pytest.main()
