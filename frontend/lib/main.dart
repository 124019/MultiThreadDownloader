import 'package:flutter/material.dart';

import 'dart:async';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    final lightScheme = ColorScheme.fromSeed(seedColor: Colors.deepPurple);

    final darkScheme = ColorScheme.fromSeed(
      seedColor: Colors.deepPurple,
      brightness: Brightness.dark,
    );

    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(colorScheme: lightScheme, useMaterial3: true),
      darkTheme: ThemeData(
        brightness: Brightness.dark,
        colorScheme: darkScheme,
        useMaterial3: true,
        scaffoldBackgroundColor: darkScheme.surface,
        appBarTheme: AppBarTheme(
          backgroundColor: darkScheme.surface,
          foregroundColor: darkScheme.onSurface,
        ),
      ),
      themeMode: ThemeMode.system, // Follow system theme mode
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

  static const String _defaultMidStr = 'You have pushed the button this many times:';
  String _midStr = _defaultMidStr;

  void _updateMessage() {
    _midStr = _counter > 10
        ? 'You have pushed the button too many times!'
        : _defaultMidStr;
  }

  void _incrementCounter() {
    setState(() {
      _counter++;
      _updateMessage();
    });
  }

  void _decrementCounter() {
    setState(() {
      _counter--;
      _updateMessage();
    });
  }

  void _onMidButtonClicked() {
    _renameTimer?.cancel(); // Cancel any existing timer and schedule a new one
    const url = 'https://go.dev';
    _incrementCounter();

    setState(() {
      _newButtonName = url;
    });

    _renameTimer = Timer(const Duration(seconds: 2), () {
      if (!mounted) return;
      setState(() {
        _newButtonName = 'Official Website of Go';
      });
    });
  }

  Timer? _renameTimer;
  @override
  void dispose() {
    _renameTimer?.cancel();
    super.dispose();
  }

  void _resetAll() {
    _renameTimer?.cancel();
    setState(() {
      _counter = 0;
      _midStr = _defaultMidStr;
    });
  }

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return Scaffold(
      appBar: AppBar(
        backgroundColor: colorScheme.surface,
        foregroundColor: colorScheme.onSurface,
        title: Text(widget.title),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(_midStr),
            Text(
              '$_counter',
              style: Theme.of(context).textTheme.headlineMedium,
            ),
            const SizedBox(height: 30),
            ElevatedButton(
              onPressed: _onMidButtonClicked,
              style: ButtonStyle(
                backgroundColor: WidgetStateProperty.resolveWith<Color>((
                  states,
                ) {
                  if (states.contains(WidgetState.hovered)) {
                    return colorScheme.primaryContainer;
                  }
                  return colorScheme.primary;
                }),
                foregroundColor: WidgetStateProperty.resolveWith<Color>((
                  states,
                ) {
                  if (states.contains(WidgetState.hovered)) {
                    return colorScheme.onPrimaryContainer;
                  }
                  return colorScheme.onPrimary;
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
            onPressed: _resetAll,
            tooltip: 'Reset',
            child: const Icon(Icons.restart_alt),
          ),
        ],
      ),
    );
  }
}
