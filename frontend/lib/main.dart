import 'package:flutter/material.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        colorScheme: .fromSeed(seedColor: Colors.deepPurple),
      ),
      home: const MyHomePage(title: 'Application Home Page'),
    );
  }
}

class MyHomePage extends StatefulWidget {
  const MyHomePage({super.key, required this.title});

  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  int _counter = 0;
  String _newButtonName = 'Official Website of Go';
  String _midStr = 'You have pushed the button this many times:';
  String midStr = 'You have pushed the button this many times:';

  void _incrementCounter() {
    setState(() {
      if (_counter > 9) {
        _midStr = 'You have pushed the button too many times!';
      }
      _counter++;
    });
  }

  void _decrementCounter() {
    setState(() {
      _counter--;
      if (_counter < 9) {
        _midStr = midStr;
      }
    });
  }

  Future<void> _fetchURL() async {
    const url = 'https://go.dev';
    setState(() {
      _midStr = midStr;
    });
    setState(() {
      _newButtonName = url;
    });
    Future.delayed(Duration(seconds: 2), () {
      setState(() {
        _newButtonName = 'Official Website of Go';
      });
    });
    _incrementCounter();
  }

  @override
  Widget build(BuildContext context) {

    return Scaffold(
      appBar: AppBar(
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        title: Text(widget.title),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: .center,
          children: [
            Text(_midStr),
            Text(
              '$_counter',
              style: Theme.of(context).textTheme.headlineMedium,
            ),
            const SizedBox(height: 30),
            ElevatedButton(
              onPressed: _fetchURL,
              style: ButtonStyle(
                backgroundColor: WidgetStateProperty.resolveWith<Color>((
                  states,
                ) {
                  if (states.contains(WidgetState.hovered)) {
                    return Colors.green;
                  }
                  return Colors.blue;
                }),
                foregroundColor: WidgetStateProperty.resolveWith<Color>((
                  states,
                ) {
                  if (states.contains(WidgetState.hovered)) {
                    return Colors.white;
                  }
                  return Colors.white;
                }),
              ),
              child: Text(_newButtonName),
            ),
          ],
        ),
      ),
      floatingActionButton: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          FloatingActionButton(
            heroTag: 'decrementBtn',
            onPressed: _decrementCounter,
            tooltip: 'Decrement',
            child: const Icon(Icons.remove),
          ),
          const SizedBox(width: 12),
          FloatingActionButton(
            heroTag: 'incrementBtn',
            onPressed: _incrementCounter,
            tooltip: 'Increment',
            child: const Icon(Icons.add),
          ),
          const SizedBox(width: 12),
          FloatingActionButton(
            heroTag: 'resetBtn',
            onPressed: () {
              setState(() {
                _counter = 0;
                _midStr = 'You have pushed the button this many times:';
              });
            },
            tooltip: 'Reset',
            child: const Icon(Icons.restart_alt),
          ),
        ],
      ),
    );
  }
}
