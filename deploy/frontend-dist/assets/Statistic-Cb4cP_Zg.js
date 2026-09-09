import{Cr as e,Yi as t,_i as n,_r as r,ir as i,na as a,ra as o,tr as s,vi as c,vr as l}from"./router-BiT_07AU.js";import{n as u}from"./index-D9IIxAgS.js";var d=n(`statistic`,[c(`label`,`
 font-weight: var(--n-label-font-weight);
 transition: .3s color var(--n-bezier);
 font-size: var(--n-label-font-size);
 color: var(--n-label-text-color);
 `),n(`statistic-value`,`
 margin-top: 4px;
 font-weight: var(--n-value-font-weight);
 `,[c(`prefix`,`
 margin: 0 4px 0 0;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-prefix-text-color);
 `,[n(`icon`,{verticalAlign:`-0.125em`})]),c(`content`,`
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-text-color);
 `),c(`suffix`,`
 margin: 0 0 0 4px;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-suffix-text-color);
 `,[n(`icon`,{verticalAlign:`-0.125em`})])])]),f=a({name:`Statistic`,props:Object.assign(Object.assign({},s.props),{tabularNums:Boolean,label:String,value:[String,Number]}),slots:Object,setup(e){let{mergedClsPrefixRef:n,inlineThemeDisabled:a,mergedRtlRef:o}=l(e),c=s(`Statistic`,`-statistic`,d,u,e,n),f=i(`Statistic`,o,n),p=t(()=>{let{self:{labelFontWeight:e,valueFontSize:t,valueFontWeight:n,valuePrefixTextColor:r,labelTextColor:i,valueSuffixTextColor:a,valueTextColor:o,labelFontSize:s},common:{cubicBezierEaseInOut:l}}=c.value;return{"--n-bezier":l,"--n-label-font-size":s,"--n-label-font-weight":e,"--n-label-text-color":i,"--n-value-font-weight":n,"--n-value-font-size":t,"--n-value-prefix-text-color":r,"--n-value-suffix-text-color":a,"--n-value-text-color":o}}),m=a?r(`statistic`,void 0,p,e):void 0;return{rtlEnabled:f,mergedClsPrefix:n,cssVars:a?void 0:p,themeClass:m?.themeClass,onRender:m?.onRender}},render(){var t;let{mergedClsPrefix:n,$slots:{default:r,label:i,prefix:a,suffix:s}}=this;return(t=this.onRender)==null||t.call(this),o(`div`,{class:[`${n}-statistic`,this.themeClass,this.rtlEnabled&&`${n}-statistic--rtl`],style:this.cssVars},e(i,e=>o(`div`,{class:`${n}-statistic__label`},this.label||e)),o(`div`,{class:`${n}-statistic-value`,style:{fontVariantNumeric:this.tabularNums?`tabular-nums`:``}},e(a,e=>e&&o(`span`,{class:`${n}-statistic-value__prefix`},e)),this.value===void 0?e(r,e=>e&&o(`span`,{class:`${n}-statistic-value__content`},e)):o(`span`,{class:`${n}-statistic-value__content`},this.value),e(s,e=>e&&o(`span`,{class:`${n}-statistic-value__suffix`},e))))}});export{f as t};