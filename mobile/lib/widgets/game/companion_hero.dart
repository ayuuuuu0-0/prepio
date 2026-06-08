import 'package:flutter/material.dart';
import '../../core/theme/design_tokens.dart';

/// CompanionHero shows a large animated companion avatar.
class CompanionHero extends StatefulWidget {
  const CompanionHero({super.key, this.name, this.species, this.size = 96});

  final String? name;
  final String? species;
  final double size;

  @override
  State<CompanionHero> createState() => _CompanionHeroState();
}

class _CompanionHeroState extends State<CompanionHero> with SingleTickerProviderStateMixin {
  late final AnimationController _controller;
  late final Animation<double> _bounce;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(vsync: this, duration: const Duration(milliseconds: 2000))
      ..repeat(reverse: true);
    _bounce = Tween<double>(begin: 0, end: -8).animate(CurvedAnimation(parent: _controller, curve: Curves.easeInOut));
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final visual = companionFor(name: widget.name, species: widget.species);

    return AnimatedBuilder(
      animation: _bounce,
      builder: (context, child) => Transform.translate(offset: Offset(0, _bounce.value), child: child),
      child: Container(
        width: widget.size,
        height: widget.size,
        decoration: BoxDecoration(
          shape: BoxShape.circle,
          gradient: LinearGradient(colors: visual.colors),
          boxShadow: [BoxShadow(color: visual.colors.first.withValues(alpha: 0.4), blurRadius: 16, offset: const Offset(0, 6))],
        ),
        alignment: Alignment.center,
        child: Text(visual.emoji, style: TextStyle(fontSize: widget.size * 0.45)),
      ),
    );
  }
}
